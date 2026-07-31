package controller

import (
	"backend/gateway/internal/client/rpc"
	"backend/gateway/internal/client/rpc/core-rpc/articlepb"
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
	var req articlepb.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	resp, err := ct.rpc.ArticleClient.CreateArticle(c, &req)
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	reponse.Success(c, resp)
}

func (ct *ArticleController) GetDetail(c *gin.Context) {
	var req articlepb.GetArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	resp, err := ct.rpc.ArticleClient.GetArticle(c, &req)
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	reponse.Success(c, resp)
}

func (ct *ArticleController) List(c *gin.Context) {
	var req articlepb.ListArticlesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	resp, err := ct.rpc.ArticleClient.ListArticles(c, &req)
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	reponse.Success(c, resp)
}

func (ct *ArticleController) ListByUserID(c *gin.Context) {
	var req articlepb.ListMyArticlesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	resp, err := ct.rpc.ArticleClient.ListMyArticles(c, &req)
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	reponse.Success(c, resp)
}

func (ct *ArticleController) ByCategory(c *gin.Context) {
	var req articlepb.ListByCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	resp, err := ct.rpc.ArticleClient.ListByCategory(c, &req)
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	reponse.Success(c, resp)
}

func (ct *ArticleController) Search(c *gin.Context) {
	var req articlepb.SearchArticlesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	resp, err := ct.rpc.ArticleClient.SearchArticles(c, &req)
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	reponse.Success(c, resp)
}

func (ct *ArticleController) Delete(c *gin.Context) {
	var req articlepb.DeleteArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	resp, err := ct.rpc.ArticleClient.DeleteArticle(c, &req)
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	reponse.Success(c, resp)
}
