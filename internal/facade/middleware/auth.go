package middleware

import (
	"errors"
	"gateway/internal/config"
	"gateway/internal/model/reponse"
	"gateway/internal/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AuthUserIDContextKey       = "user_id"
	AuthRoleContextKey         = "role"
	refreshedAccessTokenHeader = "Authorization"
)

func RequireAuth(cfg config.AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, ok := bearerToken(c.GetHeader("Authorization"))
		if !ok {
			c.Abort()
			reponse.Fail(c, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		claims, err := utils.GetClaims(accessToken, cfg.JWTSigningKey())
		if err == nil {
			if !validClaims(claims, utils.TokenTypeAccess) {
				c.Abort()
				reponse.Fail(c, http.StatusUnauthorized, "invalid or expired token")
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

		refreshClaims, err := utils.GetClaims(refreshToken, cfg.JWTSigningKey())
		if err != nil || !validClaims(refreshClaims, utils.TokenTypeRefresh) {
			utils.ClearRefreshCookie(c.Writer, cfg)
			c.Abort()
			reponse.Fail(c, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		newAccessToken, err := utils.CreateAccessToken(cfg.JWTSigningKey(), *refreshClaims, cfg.AccessDuration())
		if err != nil {
			log.Printf("auth middleware create access token: %v", err)
			c.Abort()
			reponse.Fail(c, http.StatusInternalServerError, "internal server error")
			return
		}

		c.Header(refreshedAccessTokenHeader, "Bearer "+newAccessToken)
		c.Header("Access-Control-Expose-Headers", refreshedAccessTokenHeader)
		c.Set(AuthUserIDContextKey, claims.UserID)
		c.Set(AuthRoleContextKey, claims.Role)
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
