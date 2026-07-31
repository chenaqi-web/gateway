package router

import (
	"backend/gateway/internal/facade/controller"

	"github.com/gin-gonic/gin"
)

func NewCategoryRouter(v *gin.RouterGroup, ct *controller.CategoryController) {
	category := v.Group("/types")
	{
		category.POST("/create", ct.CreateType)
		category.DELETE("/del", ct.DeleteType)
		category.GET("/list", ct.ListTypes)
		category.POST("/category/create", ct.CreateCategory)
		category.DELETE("/category/del", ct.DeleteCategory)
		category.POST("/category/list", ct.ListCategories)
	}
}
