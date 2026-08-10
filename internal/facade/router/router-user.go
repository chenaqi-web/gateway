package router

import (
	"gateway/internal/facade/controller"

	"github.com/gin-gonic/gin"
)

func NewUserRouter(v *gin.RouterGroup, userCtrl *controller.UserController, authMiddleware gin.HandlerFunc) {
	user := v.Group("/user")
	user.Use(authMiddleware)
	{
		user.GET("/", userCtrl.Get)
		user.GET("/profile", userCtrl.GetProfile)
		user.PUT("/profile", userCtrl.UpdateProfile)
	}
}
