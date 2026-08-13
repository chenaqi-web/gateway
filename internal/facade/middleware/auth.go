package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"gateway/internal/config"
	"gateway/internal/infras/cache"
	"gateway/internal/model/reponse"
	"gateway/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AuthUserIDContextKey       = "user_id"
	AuthRoleContextKey         = "role"
	refreshedAccessTokenHeader = "Authorization"
)

type AuthMiddleware struct {
	cfg          config.AuthConfig
	jwtBlackList *cache.JwtBlacklist
}

func NewAuthMiddleware(cfg config.AuthConfig, jwtBlackList *cache.JwtBlacklist) *AuthMiddleware {
	return &AuthMiddleware{
		cfg:          cfg,
		jwtBlackList: jwtBlackList,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	cfg := m.cfg
	jwtBlackList := m.jwtBlackList
	return func(c *gin.Context) {
		accessToken, ok := bearerToken(c.GetHeader("Authorization"))
		if !ok {
			c.Abort()
			reponse.Fail(c, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		blacklisted, err := jwtBlackList.IsTokenBlacklisted(c.Request.Context(), accessToken)
		if err != nil {
			log.Printf("auth middleware check access blacklist: %v", err)
			c.Abort()
			reponse.Fail(c, http.StatusInternalServerError, "internal server error")
			return
		}
		if blacklisted {
			utils.ClearRefreshCookie(c.Writer, cfg)
			c.Abort()
			reponse.Fail(c, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		// 校验access token
		claims, err := utils.GetClaims(accessToken, []byte(cfg.JWTSecret))
		if err == nil {
			if !validClaims(claims, utils.TokenTypeAccess) {
				c.Abort()
				reponse.Fail(c, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			// 校验是否拉黑
			if m.rejectBlacklistedUser(c, claims.UserID) {
				return
			}

			c.Set(AuthUserIDContextKey, claims.UserID)
			c.Set(AuthRoleContextKey, claims.Role)
			c.Next()
			return
		}
		if !isExpiredToken(err) {
			c.Abort()
			reponse.Fail(c, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		refreshToken, err := utils.RefreshTokenFromCookie(c.Request)
		if err != nil {
			utils.ClearRefreshCookie(c.Writer, cfg)
			c.Abort()
			reponse.Fail(c, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		blacklisted, err = jwtBlackList.IsTokenBlacklisted(c.Request.Context(), refreshToken)
		if err != nil {
			log.Printf("auth middleware check refresh blacklist: %v", err)
			c.Abort()
			reponse.Fail(c, http.StatusInternalServerError, "internal server error")
			return
		}
		if blacklisted {
			utils.ClearRefreshCookie(c.Writer, cfg)
			c.Abort()
			reponse.Fail(c, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		// 使用refresh刷新access
		refreshClaims, err := utils.GetClaims(refreshToken, []byte(cfg.JWTSecret))
		if err != nil || !validClaims(refreshClaims, utils.TokenTypeRefresh) {
			utils.ClearRefreshCookie(c.Writer, cfg)
			c.Abort()
			reponse.Fail(c, http.StatusUnauthorized, "invalid or expired token")
			return
		}
		if m.rejectBlacklistedUser(c, refreshClaims.UserID) {
			return
		}

		newAccessToken, err := utils.CreateAccessToken([]byte(cfg.JWTSecret), *refreshClaims, cfg.AccessExpire)
		if err != nil {
			log.Printf("auth middleware create access token: %v", err)
			c.Abort()
			reponse.Fail(c, http.StatusInternalServerError, "internal server error")
			return
		}

		c.Header(refreshedAccessTokenHeader, "Bearer "+newAccessToken)
		c.Header("Access-Control-Expose-Headers", refreshedAccessTokenHeader)
		c.Set(AuthUserIDContextKey, refreshClaims.UserID)
		c.Set(AuthRoleContextKey, refreshClaims.Role)
		c.Next()
	}
}

func (m *AuthMiddleware) rejectBlacklistedUser(c *gin.Context, userID uint64) bool {
	blacklisted, err := m.jwtBlackList.IsBlacklisted(c.Request.Context(), userID)
	if err != nil {
		log.Printf("auth middleware check user blacklist: %v", err)
		c.Abort()
		reponse.Fail(c, http.StatusInternalServerError, "internal server error")
		return true
	}
	if !blacklisted {
		return false
	}
	utils.ClearRefreshCookie(c.Writer, m.cfg)
	c.Abort()
	reponse.Fail(c, http.StatusUnauthorized, "user is blacklisted")
	return true
}

// OptionalAuth keeps public resources available to anonymous users while
// allowing controllers to read the user identity when a valid access token
// is present.
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.TrimSpace(c.GetHeader("Authorization")) == "" {
			c.Next()
			return
		}
		m.RequireAuth()(c)
	}
}

func bearerToken(value string) (string, bool) {
	parts := strings.Fields(value)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", false
	}
	return parts[1], true
}

func isExpiredToken(err error) bool {
	return errors.Is(err, jwt.ErrTokenExpired) &&
		!errors.Is(err, jwt.ErrTokenMalformed) &&
		!errors.Is(err, jwt.ErrTokenUnverifiable) &&
		!errors.Is(err, jwt.ErrTokenSignatureInvalid)
}

func validClaims(claims *utils.JWTClaims, tokenType string) bool {
	return claims != nil && claims.TokenType == tokenType && claims.UserID > 0 && claims.Role != ""
}
