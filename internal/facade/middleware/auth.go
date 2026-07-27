package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"backend/gateway/internal/client/rpc"
	"backend/gateway/internal/client/rpc/core-rpc/authpb"
	"backend/gateway/internal/model/reponse"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func NewAuthMiddleware(rpcClient *rpc.Client) (*AuthMiddleware, error) {
	if rpcClient == nil {
		return nil, fmt.Errorf("auth middleware RPC client is nil")
	}
	return newAuthMiddleware(rpcClient.GetAuthClient(), rpcClient.GetRequestTimeout())
}

func newAuthMiddleware(client accessRPC, timeout time.Duration) (*AuthMiddleware, error) {
	if client == nil || timeout <= 0 {
		return nil, fmt.Errorf("auth middleware dependency is invalid")
	}
	return &AuthMiddleware{client: client, requestTimeout: timeout}, nil
}

func (m *AuthMiddleware) RequireAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, ok := bearerToken(c.GetHeader("Authorization"))
		if !ok {
			abortInvalidAccess(c)
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), m.requestTimeout)
		defer cancel()
		identity, err := m.client.ValidateAccess(ctx, &authpb.ValidateAccessRequest{AccessToken: accessToken})
		if err != nil {
			abortAccessRPCError(c, err)
			return
		}
		if identity.GetUserId() == 0 || identity.GetSessionId() == "" || identity.GetRole() == "" {
			abortInvalidAccess(c)
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

func abortAccessRPCError(c *gin.Context, err error) {
	switch status.Code(err) {
	case codes.PermissionDenied:
		c.AbortWithStatusJSON(http.StatusForbidden, reponse.Error(10012, "用户已禁用或权限不足"))
	case codes.Unavailable, codes.DeadlineExceeded:
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, reponse.Error(http.StatusServiceUnavailable, "鉴权服务暂不可用"))
	default:
		abortInvalidAccess(c)
	}
}

func abortInvalidAccess(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, reponse.Error(10006, "Access Token 无效或过期"))
}
