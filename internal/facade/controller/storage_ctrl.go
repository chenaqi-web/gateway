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
		reponse.Fail(c, http.StatusBadRequest, "file is required")
		return
	}

	result, err := ct.svc.Upload(c.Request.Context(), file)
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	reponse.Success(c, dto.UploadResponse{
		URL:      result.URL,
		Key:      result.Key,
		Provider: ct.svc.Provider(),
	})
}

func (ct *StorageController) Delete(c *gin.Context) {
	var input dto.DeleteUploadRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	if err := ct.svc.Delete(c.Request.Context(), input.Key); err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	reponse.Success(c, nil)
}
