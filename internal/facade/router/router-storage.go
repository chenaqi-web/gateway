package router

import (
	"backend/gateway/internal/facade/controller"

	"github.com/gin-gonic/gin"
)

func NewStorageRouter(v *gin.RouterGroup, storageCtrl *controller.StorageController) {
	storage := v.Group("/storage")
	{
		storage.POST("/upload", storageCtrl.Upload)
		storage.POST("/delete", storageCtrl.Delete)
	}
}
