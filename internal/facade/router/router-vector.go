package router

import (
	"gateway/internal/facade/controller"

	"github.com/gin-gonic/gin"
)

func NewVectorRouter(v *gin.RouterGroup, vector *controller.VectorController, authMiddleware gin.HandlerFunc) {
	group := v.Group("/vector")
	group.Use(authMiddleware)
	{
		group.GET("/collections", vector.ListCollections)
		group.POST("/collections/:name", vector.CreateCollection)
		group.DELETE("/collections/:name", vector.DeleteCollection)
		group.GET("/collections/:name/documents", vector.ListDocuments)
	}
}
