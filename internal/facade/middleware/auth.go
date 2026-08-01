package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"backend/gateway/internal/client/rpc"
	"backend/gateway/internal/client/rpc/core-rpc/authpb"
	"backend/gateway/internal/model/reponse"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	AuthUserIDContextKey    = "user_id"
	AuthSessionIDContextKey = "session_id"
	AuthRoleContextKey      = "role"
)

func RequireAuth(rpcClient *rpc.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, ok := bearerToken(c.GetHeader("Authorization"))
		if !ok {
			c.Abort()
			reponse.Fail(c, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), rpcClient.GetRequestTimeout())
		defer cancel()

		identity, err := rpcClient.GetAuthClient().ValidateAccess(ctx, &authpb.ValidateAccessRequest{
			AccessToken: accessToken,
		})
		if err != nil {
			log.Printf("auth middleware: %v", err)
			c.Abort()
			switch status.Code(err) {
			case codes.Unauthenticated:
				reponse.Fail(c, http.StatusUnauthorized, "invalid or expired token")
			case codes.PermissionDenied:
				reponse.Fail(c, http.StatusForbidden, "permission denied")
			default:
				reponse.Fail(c, http.StatusBadGateway, "core-server unavailable")
			}
			return
		}
		if identity == nil || identity.GetUserId() == 0 || identity.GetSessionId() == "" || identity.GetRole() == "" {
			c.Abort()
			reponse.Fail(c, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		c.Set(AuthUserIDContextKey, identity.GetUserId())
		c.Set(AuthSessionIDContextKey, identity.GetSessionId())
		c.Set(AuthRoleContextKey, identity.GetRole())
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
