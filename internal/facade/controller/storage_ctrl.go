package controller

import (
	"backend/gateway/internal/application"
	"backend/gateway/internal/model/dto"
	"backend/gateway/internal/model/reponse"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StorageController struct {
	svc *application.StorageService
}

func NewStorageController(svc *application.StorageService) *StorageController {
	return &StorageController{svc: svc}
}

func (ct *StorageController) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, reponse.Error(http.StatusBadRequest, "file is required"))
		return
	}

	result, err := ct.svc.Upload(c.Request.Context(), file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, reponse.Success(dto.UploadResponse{
		URL:      result.URL,
		Key:      result.Key,
		Provider: ct.svc.Provider(),
	}))
}

func (ct *StorageController) Delete(c *gin.Context) {
	var input dto.DeleteUploadRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, reponse.Error(http.StatusBadRequest, "invalid request parameters"))
		return
	}

	if err := ct.svc.Delete(c.Request.Context(), input.Key); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, reponse.Success(nil))
}
