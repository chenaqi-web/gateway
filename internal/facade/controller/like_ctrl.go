package controller

import (
	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/likepb"
	"gateway/internal/facade/middleware"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LikeController struct {
	rpc *rpc.Client
}

func NewLikeController(rpcClient *rpc.Client) *LikeController {
	return &LikeController{rpc: rpcClient}
}

func (ct *LikeController) ThumbUp(c *gin.Context) {
	var req dto.LikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	resp, err := ct.rpc.LikeClient.ThumbUp(c, &likepb.ThumbUpRequest{
		UserID:     userID,
		ObjectType: req.ObjectType,
		ObjectID:   req.ObjectID,
	})
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	reponse.Success(c, dto.LikeBoolResponse{Success: resp.GetSuccess()})
}

func (ct *LikeController) CancelThumbUp(c *gin.Context) {
	var req dto.LikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	resp, err := ct.rpc.LikeClient.CancelThumbUp(c, &likepb.CancelThumbUpRequest{
		UserID:     userID,
		ObjectType: req.ObjectType,
		ObjectID:   req.ObjectID,
	})
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	reponse.Success(c, dto.LikeBoolResponse{Success: resp.GetSuccess()})
}

func (ct *LikeController) HasLike(c *gin.Context) {
	var req dto.LikeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}

	resp, err := ct.rpc.LikeClient.HasLike(c, &likepb.HasArticleLikeRequest{
		UserID:     userID,
		ObjectType: req.ObjectType,
		ObjectID:   req.ObjectID,
	})
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	reponse.Success(c, dto.LikeStatus{ObjectID: req.ObjectID, IsLiked: resp.GetIsLiked()})
}

func (ct *LikeController) BatchLikeStatus(c *gin.Context) {
	var req dto.BatchLikeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}

	resp, err := ct.rpc.LikeClient.BatchLikeStatus(c, &likepb.BatchCommentLikeStatusRequest{
		UserID:     userID,
		ObjectType: req.ObjectType,
		ObjectIDs:  req.ObjectIDs,
	})
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	items := make([]*dto.LikeStatus, 0, len(resp.GetItems()))
	for _, item := range resp.GetItems() {
		if item != nil {
			items = append(items, &dto.LikeStatus{ObjectID: item.GetObjectID(), IsLiked: item.GetIsLiked()})
		}
	}
	reponse.Success(c, dto.BatchLikeStatusResponse{Items: items})
}

func (ct *LikeController) UserLikeList(c *gin.Context) {
	var req dto.UserLikeListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}

	resp, err := ct.rpc.LikeClient.PageQueryUserLikeList(c, &likepb.PageQueryUserLikeListRequest{
		UserID:     userID,
		ObjectType: req.ObjectType,
		Page:       req.Page,
		PageSize:   req.PageSize,
	})
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	reponse.Success(c, resp)
}
