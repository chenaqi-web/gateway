package controller

import (
	"gateway/internal/application"
	"gateway/internal/facade/middleware"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"

	"github.com/gin-gonic/gin"
)

type LikeController struct{ svc *application.LikeService }

func NewLikeController(svc *application.LikeService) *LikeController {
	return &LikeController{svc: svc}
}

func (ct *LikeController) ThumbUp(c *gin.Context) {
	var req dto.LikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Unauthorized(c)
		return
	}
	req.UserID = userID
	result, err := ct.svc.ThumbUp(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c)
		return
	}
	reponse.Success(c, result)
}

func (ct *LikeController) CancelThumbUp(c *gin.Context) {
	var req dto.LikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Unauthorized(c)
		return
	}
	req.UserID = userID
	result, err := ct.svc.CancelThumbUp(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c)
		return
	}
	reponse.Success(c, result)
}

func (ct *LikeController) UserLikeList(c *gin.Context) {
	var req dto.UserLikeListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.StatusBadRequest(c)
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Unauthorized(c)
		return
	}
	req.UserID = userID
	result, err := ct.svc.UserLikeList(c.Request.Context(), req)
	if err != nil {
		reponse.InternalServerError(c)
		return
	}
	reponse.Success(c, result)
}
