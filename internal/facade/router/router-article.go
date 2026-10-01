package router

import (
	"gateway/internal/facade/controller"

	"github.com/gin-gonic/gin"
)

func NewArticleRouter(v *gin.RouterGroup, ct *controller.ArticleController, authMiddleware gin.HandlerFunc, optionalAuth gin.HandlerFunc) {
	article := v.Group("/article")
	{
		article.POST("/message", optionalAuth, ct.GetDetail)
		article.POST("/search", ct.Search)

		article.POST("/list", ct.List)
		article.POST("/list/by_cate", ct.ListByCate)
		article.POST("/list/by_user", ct.ListByUser)

		authorized := article.Group("")
		authorized.Use(authMiddleware)
		{
			// 文章操作
			authorized.POST("/create", ct.Create)
			authorized.POST("/edit", ct.Edit)
			authorized.DELETE("/del", ct.Delete)

			// 草稿箱操作
			authorized.POST("/draft", ct.SaveDraft)
			authorized.POST("/draft/list", ct.ListDrafts)
			authorized.POST("/draft/publish", ct.PublishDraft)
			authorized.DELETE("/draft/del", ct.DeleteDraft)

			// 图片操作
			authorized.POST("/upload/cover", ct.UploadCover)
			authorized.POST("/upload/content", ct.UploadContentImage)

		}
	}
}
