package router

import (
	"gateway/internal/facade/controller"

	"github.com/gin-gonic/gin"
)

func NewAIRouter(v *gin.RouterGroup, aiChat *controller.AiChatController, authMiddleware gin.HandlerFunc) {
	ai := v.Group("/ai-chat")
	ai.Use(authMiddleware)
	{
		ai.POST("/session", aiChat.CreateSession)
		ai.GET("/sessions", aiChat.ListSessions)
		ai.GET("/session/:id", aiChat.GetSession)
		ai.GET("/session/:id/messages", aiChat.ListMessages)
		ai.POST("/chat", aiChat.Chat)
	}
}
