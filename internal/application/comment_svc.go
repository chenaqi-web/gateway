package application

import (
	"context"

	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/commentpb"
	"gateway/internal/client/rpc/core-rpc/likepb"
	"gateway/internal/infras/clog"
	"gateway/internal/model/dto"
	"gateway/internal/model/enum"

	"go.uber.org/zap"
)

type CommentService struct {
	rpc *rpc.Client
	log *clog.Log
}

func NewCommentService(rpcClient *rpc.Client, log *clog.Log) *CommentService {
	return &CommentService{
		rpc: rpcClient,
		log: log,
	}
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
	comments := dto.ToCommentList(resp.GetComments())
	s.markLikedComments(ctx, req.UserID, comments)
	return &dto.CommentListResponse{
		Comments: comments,
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
	replies := dto.ToCommentList(resp.GetReplies())
	s.markLikedComments(ctx, req.UserID, replies)
	return &dto.CommentRepliesResponse{
		Replies: replies,
		Page:    resp.GetPage(),
		Size:    resp.GetSize(),
	}, nil
}

// =====================================================================================================================

// todo 暂时先这样（应该改成批量查询才对），由于列表我是展示5个，所以速度应该比较快
func (s *CommentService) markLikedComments(ctx context.Context, userID uint64, comments []*dto.CommentInfo) {
	if userID == 0 {
		return
	}
	for _, comment := range comments {
		if comment == nil {
			continue
		}
		resp, err := s.rpc.LikeClient.HasLike(ctx, &likepb.HasLikeRequest{
			UserID:     userID,
			ObjectType: enum.ObjectTypeComment.String(),
			ObjectID:   comment.ID,
		})
		if err != nil {
			s.log.Error("CommentService/HasCommentLike error", zap.Error(err))
			continue
		}
		comment.IsLiked = resp.GetIsLiked()
	}
}
