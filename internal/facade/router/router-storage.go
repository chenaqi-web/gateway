package router

import (
	"gateway/internal/facade/controller"

	"github.com/gin-gonic/gin"
)

func NewStorageRouter(v *gin.RouterGroup, storageCtrl *controller.StorageController, authMiddleware gin.HandlerFunc) {
	storage := v.Group("/storage")
	storage.Use(authMiddleware)
	{
		storage.POST("/content", storageCtrl.UploadContent)
		storage.POST("/avatar", storageCtrl.UploadAvatar)
		storage.POST("/cover", storageCtrl.UploadCover)
		storage.POST("/delete", storageCtrl.Delete)
	}
}
