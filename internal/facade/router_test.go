package facade

import (
	"testing"

	"backend/gateway/internal/config"
	"backend/gateway/internal/facade/controller"

	"github.com/gin-gonic/gin"
)

func TestGatewayRoutesKeepHealthAndAddAuthEndpoints(t *testing.T) {
	engine := New(
		&config.Config{Server: config.ServerConfig{Mode: gin.TestMode}},
		&controller.HealthController{},
		&controller.AuthController{},
		&controller.AiChatController{},
		&controller.UserController{},
		&controller.StorageController{},
	)
	routes := make(map[string]struct{})
	for _, route := range engine.Routes() {
		routes[route.Method+" "+route.Path] = struct{}{}
	}
	for _, expected := range []string{
		"GET /api/v1/health/ping",
		"POST /api/v1/auth/send-email-code",
		"POST /api/v1/auth/register",
		"POST /api/v1/auth/login",
		"POST /api/v1/auth/email-login",
		"POST /api/v1/auth/refresh",
		"POST /api/v1/auth/logout",
		"POST /api/v1/auth/reset-password-by-email",
	} {
		if _, ok := routes[expected]; !ok {
			t.Fatalf("missing route %s", expected)
		}
	}
}
