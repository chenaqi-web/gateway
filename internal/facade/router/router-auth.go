package router

import (
	"gateway/internal/facade/controller"

	"github.com/gin-gonic/gin"
)

func NewAuthRouter(v *gin.RouterGroup, auth *controller.AuthController, authMiddleware gin.HandlerFunc) {
	group := v.Group("/auth")
	{
		group.POST("/send-email-code", auth.SendEmailCode)
		group.POST("/register", auth.Register)
		group.POST("/login", auth.Login)
		group.POST("/email_login", auth.EmailLogin)
		group.POST("/forgot-password", auth.ForgotPassword)

		authorized := group.Group("")
		authorized.Use(authMiddleware)
		{
			authorized.GET("/logout", authMiddleware, auth.Logout)
		}
	}
}
