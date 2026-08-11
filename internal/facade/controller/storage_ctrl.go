package controller

import (
	"context"
	"mime/multipart"
	"net/http"

	"gateway/internal/application"
	"gateway/internal/facade/middleware"
	"gateway/internal/model/dto"
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

func (ct *StorageController) UploadCover(c *gin.Context)   { ct.upload(c, ct.svc.UploadCover) }
func (ct *StorageController) UploadContent(c *gin.Context) { ct.upload(c, ct.svc.UploadContent) }

func (ct *StorageController) UploadAvatar(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		reponse.Fail(c, http.StatusBadRequest, "file is required")
		return
	}
	result, err := ct.svc.UploadAvatar(c.Request.Context(), file)
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	if _, err := ct.userSvc.UpdateAvatar(c.Request.Context(), userID, result.URL); err != nil {
		_ = ct.svc.Delete(c.Request.Context(), result.Key)
		userRPCError(c, err)
		return
	}
	reponse.Success(c, result)
}

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
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}
	if err := ct.svc.Delete(c.Request.Context(), req.Key); err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	reponse.Success(c, nil)
}
