package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"backend/gateway/internal/client/rpc/core-rpc/authpb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeAccessRPC struct {
	called   int
	token    string
	response *authpb.ValidateAccessResponse
	err      error
}

func (f *fakeAccessRPC) ValidateAccess(_ context.Context, request *authpb.ValidateAccessRequest, _ ...grpc.CallOption) (*authpb.ValidateAccessResponse, error) {
	f.called++
	f.token = request.GetAccessToken()
	return f.response, f.err
}

func TestAuthMiddlewareUsesOnlyBearerAndSetsContext(t *testing.T) {
	client := &fakeAccessRPC{response: &authpb.ValidateAccessResponse{UserId: 7, SessionId: "session", Role: "admin"}}
	middleware, err := newAuthMiddleware(client, time.Second)
	if err != nil {
		t.Fatalf("newAuthMiddleware() error = %v", err)
	}
	engine := gin.New()
	engine.Use(middleware.RequireAccess())
	engine.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"user_id":    c.GetUint64(AuthUserIDContextKey),
			"session_id": c.GetString(AuthSessionIDContextKey),
			"role":       c.GetString(AuthRoleContextKey),
		})
	})
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Bearer access-token")
	request.AddCookie(&http.Cookie{Name: "refresh_token", Value: "must-not-be-read"})
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK || client.called != 1 || client.token != "access-token" {
		t.Fatalf("status=%d body=%s called=%d token=%q", recorder.Code, recorder.Body.String(), client.called, client.token)
	}
	for _, expected := range []string{`"user_id":7`, `"session_id":"session"`, `"role":"admin"`} {
		if !strings.Contains(recorder.Body.String(), expected) {
			t.Fatalf("response missing %s: %s", expected, recorder.Body.String())
		}
	}
}

func TestAuthMiddlewareDoesNotUseRefreshCookieOrAutoRefresh(t *testing.T) {
	client := &fakeAccessRPC{}
	middleware, err := newAuthMiddleware(client, time.Second)
	if err != nil {
		t.Fatalf("newAuthMiddleware() error = %v", err)
	}
	engine := gin.New()
	engine.Use(middleware.RequireAccess())
	engine.GET("/protected", func(c *gin.Context) { c.Status(http.StatusOK) })
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.AddCookie(&http.Cookie{Name: "refresh_token", Value: "refresh-only"})
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized || client.called != 0 {
		t.Fatalf("status=%d body=%s called=%d", recorder.Code, recorder.Body.String(), client.called)
	}
}

func TestAuthMiddlewareMapsCoreErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{name: "invalid", err: status.Error(codes.Unauthenticated, "invalid"), wantStatus: http.StatusUnauthorized},
		{name: "disabled", err: status.Error(codes.PermissionDenied, "disabled"), wantStatus: http.StatusForbidden},
		{name: "unavailable", err: status.Error(codes.Unavailable, "down"), wantStatus: http.StatusServiceUnavailable},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := &fakeAccessRPC{err: test.err}
			middleware, err := newAuthMiddleware(client, time.Second)
			if err != nil {
				t.Fatalf("newAuthMiddleware() error = %v", err)
			}
			engine := gin.New()
			engine.Use(middleware.RequireAccess())
			engine.GET("/protected", func(c *gin.Context) { c.Status(http.StatusOK) })
			request := httptest.NewRequest(http.MethodGet, "/protected", nil)
			request.Header.Set("Authorization", "Bearer access")
			recorder := httptest.NewRecorder()
			engine.ServeHTTP(recorder, request)
			if recorder.Code != test.wantStatus {
				t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
			}
		})
	}
}

func TestBearerToken(t *testing.T) {
	for _, value := range []string{"", "Basic token", "Bearer", "Bearer one two"} {
		if _, ok := bearerToken(value); ok {
			t.Fatalf("bearerToken(%q) unexpectedly valid", value)
		}
	}
	if token, ok := bearerToken("bearer token"); !ok || token != "token" {
		t.Fatalf("bearerToken() = %q, %v", token, ok)
	}
}
