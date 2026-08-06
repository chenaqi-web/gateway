package controller

import (
	"gateway/internal/application"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StorageController struct{ svc *application.StorageService }

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
	reponse.Success(c, dto.UploadResponse{URL: result.URL, Key: result.Key, Provider: ct.svc.Provider()})
}

func (ct *StorageController) Delete(c *gin.Context) {
	var req dto.DeleteUploadRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := ct.svc.Delete(c.Request.Context(), req.Key); err != nil {
		rpcError(c, err)
		return
	}
	reponse.Success(c, nil)
}
