package router

import (
	"backend/gateway/internal/facade/controller"

	"github.com/gin-gonic/gin"
)

func NewAIRouter(v *gin.RouterGroup, aiChat *controller.AiChatController) {
	ai := v.Group("/ai-chat")
	{
		ai.POST("/session", aiChat.CreateSession)
		ai.GET("/sessions", aiChat.ListSessions)
		ai.GET("/session/:id", aiChat.GetSession)
		ai.GET("/session/:id/messages", aiChat.ListMessages)
		ai.POST("/chat", aiChat.Chat)
	}
}
