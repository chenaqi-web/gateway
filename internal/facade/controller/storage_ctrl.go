package controller

import (
	"context"
	"mime/multipart"
	"net/http"

	"gateway/internal/application"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"

	"github.com/gin-gonic/gin"
)

type StorageController struct{ svc *application.StorageService }

func NewStorageController(svc *application.StorageService) *StorageController {
	return &StorageController{svc: svc}
}

func (ct *StorageController) UploadAvatar(c *gin.Context)  { ct.upload(c, ct.svc.UploadAvatar) }
func (ct *StorageController) UploadCover(c *gin.Context)   { ct.upload(c, ct.svc.UploadCover) }
func (ct *StorageController) UploadContent(c *gin.Context) { ct.upload(c, ct.svc.UploadContent) }

func (ct *StorageController) upload(c *gin.Context, handler func(context.Context, *multipart.FileHeader) (*application.UploadResponse, error)) {
	file, err := c.FormFile("file")
	if err != nil {
		reponse.Fail(c, http.StatusBadRequest, "file is required")
		return
	}
	result, err := handler(c.Request.Context(), file)
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	reponse.Success(c, result)
}

func (ct *StorageController) Delete(c *gin.Context) {
	var req dto.DeleteUploadRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := ct.svc.Delete(c.Request.Context(), req.Key); err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	reponse.Success(c, nil)
}
