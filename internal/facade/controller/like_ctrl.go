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
	req.UserID = userID
	result, err := ct.svc.ThumbUp(c.Request.Context(), req)
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
	req.UserID = userID
	result, err := ct.svc.CancelThumbUp(c.Request.Context(), req)
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
	req.UserID = userID
	result, err := ct.svc.UserLikeList(c.Request.Context(), req)
	if err != nil {
		rpcError(c, err)
		return
	}
	reponse.Success(c, result)
}
