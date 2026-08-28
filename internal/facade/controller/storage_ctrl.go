package controller

import (
	"context"
	"gateway/internal/model/dto"
	"mime/multipart"

	"gateway/internal/application"
	"gateway/internal/facade/middleware"
	"gateway/internal/model/reponse"

	"github.com/gin-gonic/gin"
)

type StorageController struct {
	svc     *application.StorageService
	userSvc *application.UserService
}

func NewStorageController(svc *application.StorageService, userSvc *application.UserService) *StorageController {
	return &StorageController{svc: svc, userSvc: userSvc}
}

func (ct *StorageController) UploadCover(c *gin.Context) { ct.upload(c, ct.svc.UploadCover) }

func (ct *StorageController) UploadContent(c *gin.Context) { ct.upload(c, ct.svc.UploadContent) }

func (ct *StorageController) UploadAvatar(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Unauthorized(c)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		reponse.StatusBadRequest(c)
		return
	}
	result, err := ct.svc.UploadAvatar(c.Request.Context(), file)
	if err != nil {
		reponse.InternalServerError(c)
		return
	}
	if _, err := ct.userSvc.UpdateAvatar(c.Request.Context(), userID, result.URL); err != nil {
		_ = ct.svc.Delete(c.Request.Context(), result.Key)
		reponse.InternalServerError(c)
		return
	}
	reponse.Success(c, result)
}

func (ct *StorageController) upload(c *gin.Context, handler func(context.Context, *multipart.FileHeader) (*application.UploadResponse, error)) {
	file, err := c.FormFile("file")
	if err != nil {
		reponse.StatusBadRequest(c)
		return
	}
	result, err := handler(c.Request.Context(), file)
	if err != nil {
		reponse.InternalServerError(c)
		return
	}
	reponse.Success(c, result)
}

func (ct *StorageController) Delete(c *gin.Context) {
	var req dto.DeleteUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}
	if err := ct.svc.Delete(c.Request.Context(), req.Key); err != nil {
		reponse.InternalServerError(c)
		return
	}
	reponse.Success(c, nil)
}
