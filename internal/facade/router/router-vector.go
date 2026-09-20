package router

import (
	"gateway/internal/facade/controller"
	"gateway/internal/facade/middleware"

	"github.com/gin-gonic/gin"
)

func NewVectorRouter(v *gin.RouterGroup, vector *controller.VectorController, authMiddleware gin.HandlerFunc) {
	v.GET("/vector/collections", vector.ListCollections)

	group := v.Group("/vector")
	// 1. 状态中间件
	group.Use(authMiddleware)
	// 2.权限中间件
	group.Use(middleware.Role())
	{
		group.POST("/collections/:name", vector.CreateCollection)
		group.DELETE("/collections/:name", vector.DeleteCollection)
		group.GET("/documents/:name", vector.ListDocuments)
	}
}
