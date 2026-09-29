package application

import (
	"context"

	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/commentpb"
	"gateway/internal/infras/clog"
	"gateway/internal/model/dto"

	"go.uber.org/zap"
)

type CommentService struct {
	rpc *rpc.Client
	log *clog.Log
}

func NewCommentService(rpcClient *rpc.Client, log *clog.Log) *CommentService {
	return &CommentService{rpc: rpcClient, log: log}
}

func (s *CommentService) Create(ctx context.Context, req dto.CreateCommentRequest) (*dto.CommentBoolResponse, error) {
	resp, err := s.rpc.CommentClient.CreateComment(ctx, &commentpb.CreateCommentReq{
		ArticleId: req.ArticleID,
		UserId:    req.UserID,
		Content:   req.Content,
	})
	if err != nil {
		s.log.Error("CommentService/Create error", zap.Error(err))
		return nil, err
	}
	return dto.ToCommentBoolResponse(resp.GetSuccess()), nil
}

func (s *CommentService) CreateReply(ctx context.Context, req dto.CreateReplyRequest) (*dto.CommentBoolResponse, error) {
	resp, err := s.rpc.CommentClient.CreateReply(ctx, &commentpb.CreateReplyReq{ArticleId: req.ArticleID,
		RootId:    req.ParentID,
		UserId:    req.UserID,
		ReplyToId: req.ReplyToID,
		Content:   req.Content,
	})
	if err != nil {
		s.log.Error("CommentService/CreateReply error", zap.Error(err))
		return nil, err
	}
	return dto.ToCommentBoolResponse(resp.GetSuccess()), nil
}

func (s *CommentService) Delete(ctx context.Context, req dto.DeleteCommentRequest) (*dto.CommentBoolResponse, error) {
	resp, err := s.rpc.CommentClient.DeleteComment(ctx, &commentpb.DeleteCommentReq{
		Id:     req.ID,
		UserId: req.UserID,
	})
	if err != nil {
		s.log.Error("CommentService/Delete error", zap.Error(err))
		return nil, err
	}
	return dto.ToCommentBoolResponse(resp.GetSuccess()), nil
}

func (s *CommentService) List(ctx context.Context, req dto.GetArticleCommentsRequest) (*dto.CommentListResponse, error) {
	resp, err := s.rpc.CommentClient.GetArticleComments(ctx, &commentpb.GetArticleCommentsReq{
		ArticleId: req.ArticleID,
		Page:      req.Page,
		Size:      req.Size,
	})
	if err != nil {
		s.log.Error("CommentService/List error", zap.Error(err))
		return nil, err
	}
	return &dto.CommentListResponse{
		Comments: dto.ToCommentList(resp.GetComments()),
		Page:     resp.GetPage(),
		Size:     resp.GetSize(),
	}, nil
}

func (s *CommentService) Replies(ctx context.Context, req dto.GetCommentRepliesRequest) (*dto.CommentRepliesResponse, error) {
	resp, err := s.rpc.CommentClient.GetCommentReplies(ctx, &commentpb.GetCommentRepliesReq{
		RootId: req.ParentID,
		Page:   req.Page,
		Size:   req.Size,
	})
	if err != nil {
		s.log.Error("CommentService/Replies error", zap.Error(err))
		return nil, err
	}
	return &dto.CommentRepliesResponse{
		Replies: dto.ToCommentList(resp.GetReplies()),
		Page:    resp.GetPage(),
		Size:    resp.GetSize(),
	}, nil
}
