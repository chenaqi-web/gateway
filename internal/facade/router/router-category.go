package router

import (
	"gateway/internal/facade/controller"

	"github.com/gin-gonic/gin"
)

func NewCategoryRouter(v *gin.RouterGroup, ct *controller.CategoryController, authMiddleware gin.HandlerFunc) {
	category := v.Group("/types")
	{
		category.GET("/list", ct.ListTypes)
		category.POST("/category/list", ct.ListCategories)

		authorized := category.Group("")
		authorized.Use(authMiddleware)
		{
			authorized.POST("/create", ct.CreateType)
			authorized.DELETE("/del", ct.DeleteType)
			authorized.POST("/category/create", ct.CreateCategory)
			authorized.DELETE("/category/del", ct.DeleteCategory)
		}
	}
}
