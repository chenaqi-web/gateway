package router

import (
	"gateway/internal/facade/controller"

	"github.com/gin-gonic/gin"
)

func NewAuthRouter(v *gin.RouterGroup, auth *controller.AuthController) {
	group := v.Group("/auth")
	{
		group.POST("/send-email-code", auth.SendEmailCode)
		group.POST("/login", auth.Login)
		group.POST("/logout", auth.Logout)
	}
}
