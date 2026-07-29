package router

import (
	"backend/gateway/internal/facade/controller"

	"github.com/gin-gonic/gin"
)

func NewAuthRouter(v *gin.RouterGroup, auth *controller.AuthController) {
	group := v.Group("/auth")
	{
		group.POST("/send-email-code", auth.SendEmailCode)
		//group.POST("/register", auth.Register)
		//group.POST("/login", auth.Login)
		//group.POST("/email-login", auth.EmailLogin)
		//group.POST("/refresh", auth.Refresh)
		//group.POST("/logout", auth.Logout)
		//group.POST("/reset-password-by-email", auth.ResetPasswordByEmail)
	}
}
