package controller

import (
	"gateway/internal/application"
	"gateway/internal/facade/middleware"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommentController struct{ svc *application.CommentService }

func NewCommentController(svc *application.CommentService) *CommentController {
	return &CommentController{svc: svc}
}

func (ct *CommentController) Create(c *gin.Context) {
	var req dto.CreateCommentRequest
	if !bindJSON(c, &req) {
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	req.UserID = userID
	result, err := ct.svc.Create(c.Request.Context(), req)
	if err != nil {
		rpcError(c, err)
		return
	}
	reponse.Success(c, result)
}

func (ct *CommentController) CreateReply(c *gin.Context) {
	var req dto.CreateReplyRequest
	if !bindJSON(c, &req) {
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	req.UserID = userID
	result, err := ct.svc.CreateReply(c.Request.Context(), req)
	if err != nil {
		rpcError(c, err)
		return
	}
	reponse.Success(c, result)
}

func (ct *CommentController) Delete(c *gin.Context) {
	var req dto.DeleteCommentRequest
	if !bindJSON(c, &req) {
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	req.UserID = userID
	result, err := ct.svc.Delete(c.Request.Context(), req)
	if err != nil {
		rpcError(c, err)
		return
	}
	reponse.Success(c, result)
}

func (ct *CommentController) List(c *gin.Context) {
	var req dto.GetArticleCommentsRequest
	if !bindJSON(c, &req) {
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	req.UserID = userID
	result, err := ct.svc.List(c.Request.Context(), req)
	if err != nil {
		rpcError(c, err)
		return
	}
	reponse.Success(c, result)
}

func (ct *CommentController) Replies(c *gin.Context) {
	var req dto.GetCommentRepliesRequest
	if !bindJSON(c, &req) {
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "authentication required")
		return
	}
	req.UserID = userID
	result, err := ct.svc.Replies(c.Request.Context(), req)
	if err != nil {
		rpcError(c, err)
		return
	}
	reponse.Success(c, result)
}
