package router

import (
	"backend/gateway/internal/facade/controller"

	"github.com/gin-gonic/gin"
)

func NewArticleRouter(v *gin.RouterGroup, ct *controller.ArticleController) {
	article := v.Group("/article")
	{
		article.POST("/create", ct.Create)
		article.DELETE("/del", ct.Delete)
		article.POST("/message", ct.GetDetail)
		article.POST("/search", ct.Search)
		article.GET("/list", ct.List)
		article.POST("/list/by_user_id", ct.ListByUserID)
		article.POST("/list/by_cate", ct.ByCategory)
	}
}
