package controller

import (
	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/categorypb"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	rpc *rpc.Client
}

func NewCategoryController(rpcClient *rpc.Client) *CategoryController {
	return &CategoryController{rpc: rpcClient}
}

func (ct *CategoryController) CreateType(c *gin.Context) {
	var req dto.CreateTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	resp, err := ct.rpc.CategoryClient.CreateType(c, &categorypb.CreateTypeRequest{Name: req.Name})
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	reponse.Success(c, dto.ToCategoryBoolResponse(resp.GetSuccess()))
}

func (ct *CategoryController) DeleteType(c *gin.Context) {
	var req dto.DeleteTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	resp, err := ct.rpc.CategoryClient.DeleteType(c, &categorypb.DeleteTypeRequest{Id: req.ID})
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	reponse.Success(c, dto.ToCategoryBoolResponse(resp.GetSuccess()))
}

func (ct *CategoryController) ListTypes(c *gin.Context) {
	resp, err := ct.rpc.CategoryClient.ListTypes(c, &categorypb.ListTypesRequest{})
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	reponse.Success(c, dto.ToListTypesResponse(resp))
}

func (ct *CategoryController) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	resp, err := ct.rpc.CategoryClient.CreateCategory(c, &categorypb.CreateCategoryRequest{
		ParentID: req.ParentID,
		Name:     req.Name,
	})
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	reponse.Success(c, dto.ToCategoryBoolResponse(resp.GetSuccess()))
}

func (ct *CategoryController) DeleteCategory(c *gin.Context) {
	var req dto.DeleteCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	resp, err := ct.rpc.CategoryClient.DeleteCategory(c, &categorypb.DeleteCategoryRequest{Id: req.ID})
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	reponse.Success(c, dto.ToCategoryBoolResponse(resp.GetSuccess()))
}

func (ct *CategoryController) ListCategories(c *gin.Context) {
	var req dto.ListCategoriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}

	resp, err := ct.rpc.CategoryClient.ListCategories(c, &categorypb.ListCategoriesRequest{
		ParentID: req.ParentID,
	})
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	reponse.Success(c, dto.ToListCategoriesResponse(resp))
}
