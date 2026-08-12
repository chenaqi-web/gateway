package controller

import (
	"gateway/internal/application"
	"gateway/internal/facade/middleware"
	"gateway/internal/model/reponse"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VectorController struct{ svc *application.VectorService }

func NewVectorController(svc *application.VectorService) *VectorController {
	return &VectorController{svc: svc}
}

func (ct *VectorController) ListCollections(c *gin.Context) {
	list, err := ct.svc.ListCollections(c.Request.Context())
	if err != nil {
		reponse.Fail(c, http.StatusBadGateway, err.Error())
		return
	}
	reponse.Success(c, list)
}

func (ct *VectorController) CreateCollection(c *gin.Context) {
	if middleware.GetRole(c) != "admin" {
		reponse.Fail(c, http.StatusForbidden, "admin access required")
		return
	}
	result, err := ct.svc.CreateCollection(c.Request.Context(), c.Param("name"))
	if err != nil {
		reponse.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	reponse.Success(c, result)
}

func (ct *VectorController) DeleteCollection(c *gin.Context) {
	if middleware.GetRole(c) != "admin" {
		reponse.Fail(c, http.StatusForbidden, "admin access required")
		return
	}
	if err := ct.svc.DeleteCollection(c.Request.Context(), c.Param("name")); err != nil {
		reponse.Fail(c, http.StatusBadGateway, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

func (ct *VectorController) ListDocuments(c *gin.Context) {
	if middleware.GetRole(c) != "admin" {
		reponse.Fail(c, http.StatusForbidden, "admin access required")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "5"))
	if pageSize > 5 {
		pageSize = 5
	}
	result, err := ct.svc.ListDocuments(c.Request.Context(), c.Param("name"), page, pageSize)
	if err != nil {
		reponse.Fail(c, http.StatusBadGateway, err.Error())
		return
	}
	reponse.Success(c, result)
}
