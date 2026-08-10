package controller

import (
	"net/http"

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
	if !bindJSON(c, &req) {
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	req.AuthorID = userID
	result, err := ct.svc.Create(c.Request.Context(), req)
	if err != nil {
		rpcError(c, err)
		return
	}
	reponse.Success(c, result)
}

func (ct *ArticleController) Search(c *gin.Context) {
	var req dto.SearchArticlesRequest
	if !bindJSON(c, &req) {
		return
	}
	result, err := ct.svc.Search(c.Request.Context(), req)
	if err != nil {
		rpcError(c, err)
		return
	}
	reponse.Success(c, result)
}

func (ct *ArticleController) Delete(c *gin.Context) {
	var req dto.DeleteArticleRequest
	if !bindJSON(c, &req) {
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	if middleware.GetRole(c) == "admin" {
		req.AuthorID = 0
	} else {
		req.AuthorID = userID
	}
	result, err := ct.svc.Delete(c.Request.Context(), req)
	if err != nil {
		rpcError(c, err)
		return
	}
	reponse.Success(c, result)
}

func (ct *ArticleController) GetDetail(c *gin.Context) {
	var req dto.GetArticleRequest
	if !bindJSON(c, &req) {
		return
	}
	result, err := ct.svc.GetDetail(c.Request.Context(), req)
	if err != nil {
		rpcError(c, err)
		return
	}
	reponse.Success(c, result)
}

func (ct *ArticleController) List(c *gin.Context) {
	var req dto.ListArticlesRequest
	if !bindJSON(c, &req) {
		return
	}
	result, err := ct.svc.List(c.Request.Context(), req)
	if err != nil {
		rpcError(c, err)
		return
	}
	reponse.Success(c, result)
}

func (ct *ArticleController) ListByUserID(c *gin.Context) {
	var req dto.ListMyArticlesRequest
	if !bindJSON(c, &req) {
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	req.AuthorID = userID
	result, err := ct.svc.ListByUserID(c.Request.Context(), req)
	if err != nil {
		rpcError(c, err)
		return
	}
	reponse.Success(c, result)
}

func (ct *ArticleController) ByCategory(c *gin.Context) {
	var req dto.ListByCategoryRequest
	if !bindJSON(c, &req) {
		return
	}
	result, err := ct.svc.ListByCategory(c.Request.Context(), req)
	if err != nil {
		rpcError(c, err)
		return
	}
	reponse.Success(c, result)
}
