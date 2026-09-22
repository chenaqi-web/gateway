package router

import (
	"gateway/internal/facade/controller"
	"gateway/internal/facade/middleware"

	"github.com/gin-gonic/gin"
)

func NewUserRouter(v *gin.RouterGroup, userCtrl *controller.UserController, authMiddleware gin.HandlerFunc) {
	// 用户信息修改
	user := v.Group("/user")
	user.Use(authMiddleware)
	{
		user.GET("/profile", userCtrl.GetProfile)
		user.PUT("/profile", userCtrl.UpdateProfile)
		user.PUT("/avatar", userCtrl.UpdateAvatar)
	}

	// 管理员对用户信息的操作
	admin := user.Group("/admin")
	admin.Use(authMiddleware)
	admin.Use(middleware.Role())
	{
		user.GET("/list", userCtrl.UserList)
		user.PUT("/status", userCtrl.UpdateStatus)
		user.POST("/search", userCtrl.SearchUser)
	}
}
