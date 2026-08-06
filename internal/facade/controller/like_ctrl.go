package controller

import (
	"gateway/internal/application"
	"gateway/internal/facade/middleware"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LikeController struct{ svc *application.LikeService }

func NewLikeController(svc *application.LikeService) *LikeController {
	return &LikeController{svc: svc}
}

func (ct *LikeController) ThumbUp(c *gin.Context) {
	var req dto.LikeRequest
	if !bindJSON(c, &req) {
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	result, err := ct.svc.ThumbUp(c.Request.Context(), userID, req)
	if err != nil {
		rpcError(c, err)
		return
	}
	reponse.Success(c, result)
}

func (ct *LikeController) CancelThumbUp(c *gin.Context) {
	var req dto.LikeRequest
	if !bindJSON(c, &req) {
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	result, err := ct.svc.CancelThumbUp(c.Request.Context(), userID, req)
	if err != nil {
		rpcError(c, err)
		return
	}
	reponse.Success(c, result)
}

func (ct *LikeController) HasLike(c *gin.Context) {
	var req dto.LikeStatusRequest
	if !bindJSON(c, &req) {
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	result, err := ct.svc.HasLike(c.Request.Context(), userID, req)
	if err != nil {
		rpcError(c, err)
		return
	}
	reponse.Success(c, result)
}

func (ct *LikeController) BatchLikeStatus(c *gin.Context) {
	var req dto.BatchLikeStatusRequest
	if !bindJSON(c, &req) {
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	result, err := ct.svc.BatchLikeStatus(c.Request.Context(), userID, req)
	if err != nil {
		rpcError(c, err)
		return
	}
	reponse.Success(c, result)
}

func (ct *LikeController) UserLikeList(c *gin.Context) {
	var req dto.UserLikeListRequest
	if !bindJSON(c, &req) {
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	result, err := ct.svc.UserLikeList(c.Request.Context(), userID, req)
	if err != nil {
		rpcError(c, err)
		return
	}
	reponse.Success(c, result)
}
