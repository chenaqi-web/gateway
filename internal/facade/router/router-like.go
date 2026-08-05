package router

import (
	"gateway/internal/facade/controller"

	"github.com/gin-gonic/gin"
)

func NewLikeRouter(v *gin.RouterGroup, ct *controller.LikeController) {
	like := v.Group("/like")
	//like.Use(middleware.RequireAuth())
	{
		like.POST("/thumb_up", ct.ThumbUp)
		like.POST("/cancel_thumb_up", ct.CancelThumbUp)
		like.POST("/has_like", ct.HasLike)
		//like.POST("/batch_status", ct.BatchLikeStatus)
		like.POST("/list", ct.UserLikeList)
	}
}
