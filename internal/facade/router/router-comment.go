package router

import (
	"gateway/internal/facade/controller"

	"github.com/gin-gonic/gin"
)

func NewCommentRouter(v *gin.RouterGroup, ct *controller.CommentController) {
	comments := v.Group("/comment")
	{
		comments.POST("/create", ct.Create)
		comments.POST("/reply", ct.CreateReply)
		comments.DELETE("/delete", ct.Delete)
		comments.POST("/list", ct.List)
		comments.POST("/replies", ct.Replies)
	}
}
