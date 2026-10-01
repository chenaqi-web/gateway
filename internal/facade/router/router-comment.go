package router

import (
	"gateway/internal/facade/controller"

	"github.com/gin-gonic/gin"
)

func NewCommentRouter(v *gin.RouterGroup, ct *controller.CommentController, authMiddleware gin.HandlerFunc, optionalAuth gin.HandlerFunc) {
	comments := v.Group("/comment")
	{
		publicReads := comments.Group("")
		publicReads.Use(optionalAuth)
		publicReads.POST("/list", ct.List)
		publicReads.POST("/replies", ct.Replies)

		authorized := comments.Group("")
		authorized.Use(authMiddleware)
		{
			authorized.POST("/create", ct.Create)
			authorized.POST("/reply", ct.CreateReply)
			authorized.DELETE("/delete", ct.Delete)
		}
	}
}
