package controller

import (
	"backend/gateway/internal/client/rpc"
	"backend/gateway/internal/client/rpc/core-rpc/articlepb"
	"backend/gateway/internal/model/dto"
	"backend/gateway/internal/model/reponse"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ArticleController struct {
	rpc *rpc.Client
}

func NewArticleController(rpcClient *rpc.Client) *ArticleController {
	return &ArticleController{rpc: rpcClient}
}

func (ct *ArticleController) Create(c *gin.Context) {
	var req dto.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	resp, err := ct.rpc.ArticleClient.CreateArticle(c, &articlepb.CreateArticleRequest{
		AuthorID:   req.AuthorID,
		Title:      req.Title,
		Summary:    req.Summary,
		Content:    req.Content,
		CoverImage: req.CoverImage,
		CategoryID: req.CategoryID,
		IsTop:      req.IsTop,
	})
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	reponse.Success(c, dto.ToArticleBoolResponse(resp.GetSuccess()))
}

func (ct *ArticleController) GetDetail(c *gin.Context) {
	var req dto.GetArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	resp, err := ct.rpc.ArticleClient.GetArticle(c, &articlepb.GetArticleRequest{Id: req.ID})
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	reponse.Success(c, dto.ToGetArticleResponse(resp))
}

func (ct *ArticleController) List(c *gin.Context) {
	var req dto.ListArticlesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	resp, err := ct.rpc.ArticleClient.ListArticles(c, &articlepb.ListArticlesRequest{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	reponse.Success(c, dto.ToListArticlesResponse(resp.GetArticles()))
}

func (ct *ArticleController) ListByUserID(c *gin.Context) {
	var req dto.ListMyArticlesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	resp, err := ct.rpc.ArticleClient.ListMyArticles(c, &articlepb.ListMyArticlesRequest{
		AuthorID: req.AuthorID,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	reponse.Success(c, dto.ToListArticlesResponse(resp.GetArticles()))
}

func (ct *ArticleController) ByCategory(c *gin.Context) {
	var req dto.ListByCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	resp, err := ct.rpc.ArticleClient.ListByCategory(c, &articlepb.ListByCategoryRequest{
		CategoryID: req.CategoryID,
		Page:       req.Page,
		PageSize:   req.PageSize,
	})
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	reponse.Success(c, dto.ToListArticlesResponse(resp.GetArticles()))
}

func (ct *ArticleController) Search(c *gin.Context) {
	var req dto.SearchArticlesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	resp, err := ct.rpc.ArticleClient.SearchArticles(c, &articlepb.SearchArticlesRequest{
		Q:        req.Q,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	reponse.Success(c, dto.ToListArticlesResponse(resp.GetArticles()))
}

func (ct *ArticleController) Delete(c *gin.Context) {
	var req dto.DeleteArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	resp, err := ct.rpc.ArticleClient.DeleteArticle(c, &articlepb.DeleteArticleRequest{
		Id:       req.ID,
		AuthorID: req.AuthorID,
	})
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	reponse.Success(c, dto.ToArticleBoolResponse(resp.GetSuccess()))
}
