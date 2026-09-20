package controller

import (
	"gateway/internal/application"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"
	"gateway/internal/utils/vaildate"
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
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, list)
}

func (ct *VectorController) CreateCollection(c *gin.Context) {
	var req dto.CreateVectorCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}

	err := ct.svc.CreateCollection(c.Request.Context(), &req)
	if err != nil {
		reponse.StatusBadRequest(c)
		return
	}
	reponse.Success(c, nil)
}

func (ct *VectorController) DeleteCollection(c *gin.Context) {
	var req dto.DelVectorCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}

	if err := ct.svc.DeleteCollection(c.Request.Context(), &req); err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, nil)
}

func (ct *VectorController) ListDocuments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// name 从查询参数取，比如 /documents?name=xxx&page=1&pageSize=20
	result, err := ct.svc.ListDocuments(c.Request.Context(), &dto.ListDocumentsRequest{
		Name:     c.Param("name"),
		Page:     vaildate.Page(page),
		PageSize: vaildate.PageSize(pageSize),
	})

	if err != nil {
		reponse.InternalServerError(c, err.Error())
		return
	}
	reponse.Success(c, result)
}
