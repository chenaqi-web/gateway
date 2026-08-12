package router

import (
	"gateway/internal/facade/controller"

	"github.com/gin-gonic/gin"
)

func NewAIRouter(v *gin.RouterGroup, aiChat *controller.AiChatController, authMiddleware gin.HandlerFunc) {
	v.GET("/ai-chat/status", aiChat.GetStatus)
	ai := v.Group("/ai-chat")
	ai.Use(authMiddleware)
	{
		ai.PUT("/status", aiChat.UpdateStatus)
		ai.GET("/settings", aiChat.GetSettings)
		ai.PUT("/settings", aiChat.UpdateSettings)

		ai.POST("/session", aiChat.CreateSession)
		ai.GET("/sessions", aiChat.ListSessions)
		ai.PUT("/session/:id", aiChat.UpdateSession)
		ai.GET("/session/:id/messages", aiChat.ListMessages)
		ai.DELETE("/session/:id", aiChat.DeleteSession)

		ai.POST("/chat", aiChat.Chat)
	}
}
