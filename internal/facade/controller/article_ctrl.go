package controller

import (
	"gateway/internal/application"
	"gateway/internal/facade/middleware"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"

	"github.com/gin-gonic/gin"
)

type ArticleController struct {
	svc *application.ArticleService
}

func NewArticleController(svc *application.ArticleService) *ArticleController {
	return &ArticleController{svc: svc}
}

func (ct *ArticleController) Create(c *gin.Context) {
	var req dto.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Unauthorized(c)
		return
	}
	req.AuthorID = userID

	result, err := ct.svc.Create(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, result)
}

func (ct *ArticleController) SaveDraft(c *gin.Context) {
	var req dto.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Unauthorized(c)
		return
	}
	req.AuthorID = userID
	req.IsPublish = false

	result, err := ct.svc.Create(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, result)
}

func (ct *ArticleController) Edit(c *gin.Context) {
	var req dto.EditArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Unauthorized(c)
		return
	}
	req.AuthorID = userID

	result, err := ct.svc.Edit(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, result)
}

func (ct *ArticleController) PublishDraft(c *gin.Context) {
	var req dto.PublishDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Unauthorized(c)
		return
	}
	req.AuthorID = userID

	result, err := ct.svc.PublishDraft(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, result)
}

func (ct *ArticleController) UploadCover(c *gin.Context) {
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

	req := dto.ArticleImageUploadRequest{
		UserID: userID,
		File:   file,
	}

	result, err := ct.svc.UploadCover(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, result)
}

func (ct *ArticleController) UploadContentImage(c *gin.Context) {
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

	req := dto.ArticleImageUploadRequest{
		UserID: userID,
		File:   file,
	}

	result, err := ct.svc.UploadContentImage(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, result)
}

func (ct *ArticleController) Search(c *gin.Context) {
	var req dto.SearchArticlesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}

	result, err := ct.svc.Search(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}

	reponse.Success(c, result)
}

func (ct *ArticleController) Delete(c *gin.Context) {
	var req dto.DeleteArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Unauthorized(c)
		return
	}

	req.AuthorID = userID
	req.Role = middleware.GetRole(c)

	result, err := ct.svc.Delete(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, result)
}

func (ct *ArticleController) DeleteDraft(c *gin.Context) {
	var req dto.DeleteDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Unauthorized(c)
		return
	}
	req.AuthorID = userID

	result, err := ct.svc.DeleteDraft(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, result)
}

func (ct *ArticleController) GetDetail(c *gin.Context) {
	var req dto.GetArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}

	result, err := ct.svc.GetDetail(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}

	reponse.Success(c, result)
}

func (ct *ArticleController) List(c *gin.Context) {
	var req dto.ListArticlesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}

	result, err := ct.svc.List(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}

	reponse.Success(c, result)
}

func (ct *ArticleController) ListByUser(c *gin.Context) {
	var req dto.ListMyArticlesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}

	if req.AuthorID == 0 {
		userID, ok := middleware.GetUserID(c)
		if !ok {
			reponse.StatusBadRequest(c)
			return
		}
		req.AuthorID = userID
	}
	if req.AuthorID == 0 {
		reponse.StatusBadRequest(c)
		return
	}
	isPublished := true
	req.IsPublished = &isPublished

	result, err := ct.svc.ListByUserID(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}

	reponse.Success(c, result)
}

func (ct *ArticleController) ListDrafts(c *gin.Context) {
	var req dto.ListMyArticlesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Unauthorized(c)
		return
	}
	isPublished := false
	req.AuthorID = userID
	req.IsPublished = &isPublished

	result, err := ct.svc.ListByUserID(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}

	reponse.Success(c, result)
}

func (ct *ArticleController) ListByCate(c *gin.Context) {
	var req dto.ListByCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}

	result, err := ct.svc.ListByCategory(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}

	reponse.Success(c, result)
}
