package middleware

import (
	"context"
	"strings"
	"time"

	"backend/gateway/internal/client/rpc/core-rpc/authpb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

const (
	AuthUserIDContextKey    = "user_id"
	AuthSessionIDContextKey = "session_id"
	AuthRoleContextKey      = "role"
)

type accessRPC interface {
	ValidateAccess(context.Context, *authpb.ValidateAccessRequest, ...grpc.CallOption) (*authpb.ValidateAccessResponse, error)
}

type AuthMiddleware struct {
	client         accessRPC
	requestTimeout time.Duration
}

func (m *AuthMiddleware) RequireAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, ok := bearerToken(c.GetHeader("Authorization"))
		if !ok {
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), m.requestTimeout)
		defer cancel()
		identity, err := m.client.ValidateAccess(ctx, &authpb.ValidateAccessRequest{AccessToken: accessToken})
		if err != nil {
			return
		}
		if identity.GetUserId() == 0 || identity.GetSessionId() == "" || identity.GetRole() == "" {
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
