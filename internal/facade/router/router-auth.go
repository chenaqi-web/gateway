package router

import (
	"gateway/internal/facade/controller"

	"github.com/gin-gonic/gin"
)

func NewAuthRouter(v *gin.RouterGroup, auth *controller.AuthController, authMiddleware gin.HandlerFunc) {
	group := v.Group("/auth")
	{
		group.POST("/send-email-code", auth.SendEmailCode)
		group.POST("/email_login", auth.EmailLogin)
		group.POST("/logout", authMiddleware, auth.Logout)
	}
}
