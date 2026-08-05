package router

import (
	"gateway/internal/facade/controller"

	"github.com/gin-gonic/gin"
)

func NewCommentRouter(v *gin.RouterGroup, ct *controller.CommentController, authMiddleware gin.HandlerFunc) {
	comments := v.Group("/comment")
	{
		comments.POST("/list", ct.List)
		comments.POST("/replies", ct.Replies)

		authorized := comments.Group("")
		authorized.Use(authMiddleware)
		{
			authorized.POST("/create", ct.Create)
			authorized.POST("/reply", ct.CreateReply)
			authorized.DELETE("/delete", ct.Delete)
		}
	}
}
