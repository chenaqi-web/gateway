package router

import (
	"gateway/internal/facade/controller"

	"github.com/gin-gonic/gin"
)

func NewArticleRouter(v *gin.RouterGroup, ct *controller.ArticleController, authMiddleware gin.HandlerFunc) {
	article := v.Group("/article")
	{
		article.POST("/message", ct.GetDetail)
		article.POST("/search", ct.Search)
		article.POST("/list", ct.List)
		article.POST("/list/by_cate", ct.ByCategory)

		authorized := article.Group("")
		authorized.Use(authMiddleware)
		{
			authorized.POST("/create", ct.Create)
			authorized.DELETE("/del", ct.Delete)
			authorized.POST("/list/by_user_id", ct.ListByUserID)
		}
	}
}
