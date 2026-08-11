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
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	req.UserID = userID
	result, err := ct.svc.Create(c.Request.Context(), req)
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	reponse.Success(c, result)
}

func (ct *CommentController) CreateReply(c *gin.Context) {
	var req dto.CreateReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	req.UserID = userID
	result, err := ct.svc.CreateReply(c.Request.Context(), req)
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	reponse.Success(c, result)
}

func (ct *CommentController) Delete(c *gin.Context) {
	var req dto.DeleteCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}
	userID, ok := middleware.GetUserID(c)
	if !ok {
		reponse.Fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	req.UserID = userID
	result, err := ct.svc.Delete(c.Request.Context(), req)
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	reponse.Success(c, result)
}

func (ct *CommentController) List(c *gin.Context) {
	var req dto.GetArticleCommentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}
	if userID, ok := middleware.GetUserID(c); ok {
		req.UserID = userID
	}
	result, err := ct.svc.List(c.Request.Context(), req)
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	reponse.Success(c, result)
}

func (ct *CommentController) Replies(c *gin.Context) {
	var req dto.GetCommentRepliesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}
	if userID, ok := middleware.GetUserID(c); ok {
		req.UserID = userID
	}
	result, err := ct.svc.Replies(c.Request.Context(), req)
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	reponse.Success(c, result)
}
