package router

import "github.com/gin-gonic/gin"

type AuthHandler interface {
	SendEmailCode(*gin.Context)
	Register(*gin.Context)
	Login(*gin.Context)
	EmailLogin(*gin.Context)
	Refresh(*gin.Context)
	Logout(*gin.Context)
	ResetPasswordByEmail(*gin.Context)
}

func NewAuthRouter(v *gin.RouterGroup, auth AuthHandler) {
	group := v.Group("/auth")
	{
		group.POST("/send-email-code", auth.SendEmailCode)
		group.POST("/register", auth.Register)
		group.POST("/login", auth.Login)
		group.POST("/email-login", auth.EmailLogin)
		group.POST("/refresh", auth.Refresh)
		group.POST("/logout", auth.Logout)
		group.POST("/reset-password-by-email", auth.ResetPasswordByEmail)
	}
}
