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
	admin := v.Group("/admin")
	admin.Use(authMiddleware)
	admin.Use(middleware.Role())
	{
		admin.GET("/list", userCtrl.UserList)
		admin.GET("/search", userCtrl.SearchUser)
		admin.PUT("/status", userCtrl.UpdateStatus)
	}
}
