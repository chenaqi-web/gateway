package middleware

import (
	"net/http"
	"strings"

	"gateway/internal/config"
	"gateway/internal/infras/cache"
	"gateway/internal/model/reponse"
	"gateway/internal/utils"

	"github.com/gin-gonic/gin"
)

const (
	AuthUserIDContextKey       = "user_id"
	AuthRoleContextKey         = "role"
	refreshedAccessTokenHeader = "Authorization"
)

type AuthMiddleware struct {
	cfg       config.AuthConfig
	Blacklist *cache.Blacklist
}

func NewAuthMiddleware(cfg config.AuthConfig, Blacklist *cache.Blacklist) *AuthMiddleware {
	return &AuthMiddleware{
		cfg:       cfg,
		Blacklist: Blacklist,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1.首先获取accessToken
		accessToken, ok := bearerToken(c.GetHeader("Authorization"))
		if !ok {
			c.Abort()
			reponse.Fail(c, http.StatusUnauthorized, "invalid or expired accesstoken")
			return
		}

		// 2.判断是否在黑名单里
		blacklisted, err := m.Blacklist.IsTokenBlacklisted(c.Request.Context(), accessToken)
		if err != nil {
			c.Abort()
			reponse.Fail(c, http.StatusInternalServerError, "internal server error")
			return
		}
		if blacklisted {
			utils.ClearRefreshCookie(c.Writer, m.cfg)
			c.Abort()
			reponse.Fail(c, http.StatusUnauthorized, "token is in blacklist")
			return
		}

		// 3.校验access token
		claims, err := utils.GetClaims(accessToken, []byte(m.cfg.JWTSecret))
		UserID, Role := claims.UserID, claims.Role
		if err != nil {
			// 4. 如果有问题，则需要申请refreshToken
			refreshToken, err := utils.RefreshTokenFromCookie(c.Request)
			if err != nil {
				c.Abort()
				reponse.Fail(c, http.StatusUnauthorized, "invalid or expired refreshtoken")
				return
			}

			// 5.判断refreshToken是否在黑名单中
			blacklisted, err = m.Blacklist.IsTokenBlacklisted(c.Request.Context(), refreshToken)
			if err != nil {
				c.Abort()
				reponse.Fail(c, http.StatusInternalServerError, "internal server error")
				return
			}
			if blacklisted {
				c.Abort()
				reponse.Fail(c, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			// 6.使用refresh刷新access
			refreshClaims, err := utils.GetClaims(refreshToken, []byte(m.cfg.JWTSecret))
			if err != nil {
				c.Abort()
				reponse.Fail(c, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			newAccessToken, err := utils.CreateAccessToken([]byte(m.cfg.JWTSecret), *refreshClaims, m.cfg.AccessExpire)
			if err != nil {
				c.Abort()
				reponse.Fail(c, http.StatusInternalServerError, "internal server error")
				return
			}

			c.Header(refreshedAccessTokenHeader, "Bearer "+newAccessToken)
			c.Header("Access-Control-Expose-Headers", refreshedAccessTokenHeader)
			UserID = refreshClaims.UserID
			Role = refreshClaims.Role
		}

		// 4.校验是否拉黑
		isUserInBlacklist, err := m.Blacklist.IsUserBlacklisted(c.Request.Context(), UserID)
		if err != nil {
			c.Abort()
			reponse.Fail(c, http.StatusInternalServerError, "internal server error")
			return
		}
		if isUserInBlacklist {
			utils.ClearRefreshCookie(c.Writer, m.cfg)
			c.Abort()
			reponse.Fail(c, http.StatusUnauthorized, "user is blacklisted")
			return
		}
		c.Set(AuthUserIDContextKey, UserID)
		c.Set(AuthRoleContextKey, Role)
		c.Next()
	}
}

func bearerToken(value string) (string, bool) {
	parts := strings.Fields(value)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", false
	}
	return parts[1], true
}
