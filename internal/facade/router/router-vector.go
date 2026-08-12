package router

import (
	"gateway/internal/facade/controller"

	"github.com/gin-gonic/gin"
)

func NewVectorRouter(v *gin.RouterGroup, vector *controller.VectorController, authMiddleware gin.HandlerFunc) {
	v.GET("/vector/collections", vector.ListCollections)
	group := v.Group("/vector")
	group.Use(authMiddleware)
	{
		group.POST("/collections/:name", vector.CreateCollection)
		group.DELETE("/collections/:name", vector.DeleteCollection)
		group.GET("/collections/:name/documents", vector.ListDocuments)
	}
}
