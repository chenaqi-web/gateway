package controller

import (
	"gateway/internal/application"
	"gateway/internal/facade/middleware"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"

	"github.com/gin-gonic/gin"
)

type CategoryController struct{ svc *application.CategoryService }

func NewCategoryController(svc *application.CategoryService) *CategoryController {
	return &CategoryController{svc: svc}
}

func (ct *CategoryController) CreateType(c *gin.Context) {
	role := middleware.GetRole(c)
	if role != "admin" {
		reponse.Forbidden(c)
		return
	}

	var req dto.CreateTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}

	result, err := ct.svc.CreateType(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c)
		return
	}
	reponse.Success(c, result)
}

func (ct *CategoryController) DeleteType(c *gin.Context) {
	role := middleware.GetRole(c)
	if role != "admin" {
		reponse.Forbidden(c)
		return
	}

	var req dto.DeleteTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}

	result, err := ct.svc.DeleteType(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c)
		return
	}
	reponse.Success(c, result)
}

func (ct *CategoryController) CreateCategory(c *gin.Context) {
	role := middleware.GetRole(c)
	if role != "admin" {
		reponse.Forbidden(c)
		return
	}

	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}

	result, err := ct.svc.CreateCategory(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c)
		return
	}
	reponse.Success(c, result)
}

func (ct *CategoryController) DeleteCategory(c *gin.Context) {
	role := middleware.GetRole(c)
	if role != "admin" {
		reponse.Forbidden(c)
		return
	}

	var req dto.DeleteCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}

	result, err := ct.svc.DeleteCategory(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c)
		return
	}
	reponse.Success(c, result)
}

func (ct *CategoryController) ListTypes(c *gin.Context) {
	result, err := ct.svc.ListTypes(c.Request.Context())
	if err != nil {
		reponse.InternalServerError(c)
		return
	}
	reponse.Success(c, result)
}

func (ct *CategoryController) ListCategories(c *gin.Context) {
	var req dto.ListCategoriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}
	result, err := ct.svc.ListCategories(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c)
		return
	}
	reponse.Success(c, result)
}
