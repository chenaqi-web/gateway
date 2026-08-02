package router

import (
	"gateway/internal/facade/controller"

	"github.com/gin-gonic/gin"
)

func NewUserRouter(v *gin.RouterGroup, userCtrl *controller.UserController) {
	user := v.Group("/user")
	{
		user.GET("/", userCtrl.Get)
	}
}
