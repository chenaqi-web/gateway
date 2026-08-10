package controller

import (
	"net/http"

	"gateway/internal/application"
	"gateway/internal/facade/middleware"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"

	"github.com/gin-gonic/gin"
)

type StorageController struct{ svc *application.StorageService }

func NewStorageController(svc *application.StorageService) *StorageController {
	return &StorageController{svc: svc}
}

func (ct *StorageController) UploadAvatar(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		reponse.Fail(c, http.StatusBadRequest, "file is required")
		return
	}
	result, err := ct.svc.UploadAvatar(c.Request.Context(), userID, file)
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	reponse.Success(c, result)
}

func (ct *StorageController) UploadCover(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Unauthorized(c)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		reponse.Fail(c, http.StatusBadRequest, "file is required")
		return
	}
	result, err := ct.svc.UploadCover(c.Request.Context(), userID, c.PostForm("sessionID"), file)
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	reponse.Success(c, result)
}

func (ct *StorageController) UploadContent(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		reponse.Fail(c, http.StatusBadRequest, "file is required")
		return
	}
	result, err := ct.svc.UploadContent(c.Request.Context(), userID, c.PostForm("sessionID"), file)
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	reponse.Success(c, result)
}

func (ct *StorageController) Delete(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	var req dto.DeleteUploadRequest
	if !bindJSON(c, &req) {
		return
	}
	if err := ct.svc.DeleteOwned(c.Request.Context(), userID, req.Key); err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	reponse.Success(c, nil)
}
