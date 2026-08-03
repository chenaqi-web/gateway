package controller

import (
	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/commentpb"
	"gateway/internal/model/dto"
	"gateway/internal/model/reponse"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommentController struct{ rpc *rpc.Client }

func NewCommentController(rpcClient *rpc.Client) *CommentController {
	return &CommentController{rpc: rpcClient}
}

func (ct *CommentController) Create(c *gin.Context) {
	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}
	resp, err := ct.rpc.CommentClient.CreateComment(c, &commentpb.CreateCommentReq{ArticleId: req.ArticleID, UserId: req.UserID, Content: req.Content})
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	reponse.Success(c, dto.CommentBoolResponse{Success: resp.GetSuccess()})
}

func (ct *CommentController) CreateReply(c *gin.Context) {
	var req dto.CreateReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}
	resp, err := ct.rpc.CommentClient.CreateReply(c, &commentpb.CreateReplyReq{ArticleId: req.ArticleID, RootId: req.ParentID, UserId: req.UserID, ReplyToId: req.ReplyToID, Content: req.Content})
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	reponse.Success(c, dto.CommentBoolResponse{Success: resp.GetSuccess()})
}

func (ct *CommentController) Delete(c *gin.Context) {
	var req dto.DeleteCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}
	resp, err := ct.rpc.CommentClient.DeleteComment(c, &commentpb.DeleteCommentReq{Id: req.ID, UserId: req.UserID})
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	reponse.Success(c, dto.CommentBoolResponse{Success: resp.GetSuccess()})
}

func (ct *CommentController) List(c *gin.Context) {
	var req dto.GetArticleCommentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}
	resp, err := ct.rpc.CommentClient.GetArticleComments(c, &commentpb.GetArticleCommentsReq{ArticleId: req.ArticleID, Page: req.Page, Size: req.Size})
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	reponse.Success(c, dto.CommentListResponse{Comments: dto.ToCommentList(resp.GetComments()), Page: resp.GetPage(), Size: resp.GetSize()})
}

func (ct *CommentController) Replies(c *gin.Context) {
	var req dto.GetCommentRepliesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		reponse.Fail(c, http.StatusBadRequest, "invalid request parameters")
		return
	}
	resp, err := ct.rpc.CommentClient.GetCommentReplies(c, &commentpb.GetCommentRepliesReq{ParentId: req.ParentID, Page: req.Page, Size: req.Size})
	if err != nil {
		reponse.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	reponse.Success(c, dto.CommentRepliesResponse{Replies: dto.ToCommentList(resp.GetReplies()), Page: resp.GetPage(), Size: resp.GetSize()})
}
