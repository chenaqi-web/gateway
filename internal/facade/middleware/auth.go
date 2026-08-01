package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"backend/gateway/internal/config"
	"backend/gateway/internal/model/reponse"
	"backend/gateway/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AuthUserIDContextKey       = "user_id"
	AuthRoleContextKey         = "role"
	refreshedAccessTokenHeader = "Authorization"
)

func RequireAuth(cfg config.AuthConfig) gin.HandlerFunc {
	signingKey, configErr := cfg.JWTSigningKey()
	accessExpiresIn, durationErr := cfg.AccessDuration()
	if configErr == nil {
		configErr = durationErr
	}
	if configErr != nil {
		log.Printf("auth middleware config: %v", configErr)
	}

	return func(c *gin.Context) {
		if configErr != nil {
			c.Abort()
			reponse.Fail(c, http.StatusInternalServerError, "internal server error")
			return
		}

		accessToken, ok := bearerToken(c.GetHeader("Authorization"))
		if !ok {
			c.Abort()
			reponse.Fail(c, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		claims, err := utils.GetClaims(accessToken, signingKey)
		if err == nil {
			if !validClaims(claims, utils.TokenTypeAccess) {
				c.Abort()
				reponse.Fail(c, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			setAuthContext(c, claims)
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

		refreshClaims, err := utils.GetClaims(refreshToken, signingKey)
		if err != nil || !validClaims(refreshClaims, utils.TokenTypeRefresh) {
			utils.ClearRefreshCookie(c.Writer, cfg)
			c.Abort()
			reponse.Fail(c, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		newAccessToken, err := utils.CreateAccessToken(signingKey, *refreshClaims, accessExpiresIn)
		if err != nil {
			log.Printf("auth middleware create access token: %v", err)
			c.Abort()
			reponse.Fail(c, http.StatusInternalServerError, "internal server error")
			return
		}

		c.Header(refreshedAccessTokenHeader, "Bearer "+newAccessToken)
		c.Header("Access-Control-Expose-Headers", refreshedAccessTokenHeader)
		setAuthContext(c, refreshClaims)
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

func isExpiredToken(err error) bool {
	return errors.Is(err, jwt.ErrTokenExpired) &&
		!errors.Is(err, jwt.ErrTokenMalformed) &&
		!errors.Is(err, jwt.ErrTokenUnverifiable) &&
		!errors.Is(err, jwt.ErrTokenSignatureInvalid)
}

func validClaims(claims *utils.JWTClaims, tokenType string) bool {
	return claims != nil && claims.TokenType == tokenType && claims.UserID > 0 && claims.Role != ""
}

func setAuthContext(c *gin.Context, claims *utils.JWTClaims) {
	c.Set(AuthUserIDContextKey, claims.UserID)
	c.Set(AuthRoleContextKey, claims.Role)
}
