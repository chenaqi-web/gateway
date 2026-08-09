package router

import (
	"gateway/internal/facade/controller"

	"github.com/gin-gonic/gin"
)

func NewCommentRouter(v *gin.RouterGroup, ct *controller.CommentController, authMiddleware gin.HandlerFunc) {
	comments := v.Group("/comment")
	comments.Use(authMiddleware)
	{
		comments.POST("/list", ct.List)
		comments.POST("/replies", ct.Replies)
		comments.POST("/create", ct.Create)
		comments.POST("/reply", ct.CreateReply)
		comments.DELETE("/delete", ct.Delete)
	}
}
