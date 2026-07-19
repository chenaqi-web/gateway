package router

import (
	"backend/gateway/internal/facade/controller"
	"backend/gateway/internal/facade/middleware"

	"github.com/gin-gonic/gin"

	"backend/gateway/internal/config"
)

func New(cfg *config.Config,
	health *controller.HealthController,
	aiChat *controller.AiChatController,
	userCtrl *controller.UserController,
) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()

	// 跨域中间件
	r.Use(middleware.Cors())
	// r.Use(middleware.Recovery(), middleware.CORS())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health/ping", health.Ping)

		ai := v1.Group("/ai-chat")
		{
			ai.POST("/session", aiChat.CreateSession)
			ai.GET("/sessions", aiChat.ListSessions)
			ai.GET("/session/:id", aiChat.GetSession)
			ai.GET("/session/:id/messages", aiChat.ListMessages)
			ai.POST("/chat", aiChat.Chat)
		}

		user := v1.Group("/user")
		{
			user.GET("/", userCtrl.Get)
		}

	}

	return r
}
