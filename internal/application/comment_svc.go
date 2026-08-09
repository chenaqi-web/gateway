package application

import (
	"context"

	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/commentpb"
	"gateway/internal/client/rpc/core-rpc/likepb"
	"gateway/internal/model/dto"
)

type CommentService struct{ rpc *rpc.Client }

func NewCommentService(rpcClient *rpc.Client) *CommentService { return &CommentService{rpc: rpcClient} }

func (s *CommentService) Create(ctx context.Context, req dto.CreateCommentRequest) (*dto.CommentBoolResponse, error) {
	resp, err := s.rpc.CommentClient.CreateComment(ctx, &commentpb.CreateCommentReq{ArticleId: req.ArticleID, UserId: req.UserID, Content: req.Content})
	if err != nil {
		return nil, err
	}
	return &dto.CommentBoolResponse{Success: resp.GetSuccess()}, nil
}

func (s *CommentService) CreateReply(ctx context.Context, req dto.CreateReplyRequest) (*dto.CommentBoolResponse, error) {
	resp, err := s.rpc.CommentClient.CreateReply(ctx, &commentpb.CreateReplyReq{ArticleId: req.ArticleID, RootId: req.ParentID, UserId: req.UserID, ReplyToId: req.ReplyToID, Content: req.Content})
	if err != nil {
		return nil, err
	}
	return &dto.CommentBoolResponse{Success: resp.GetSuccess()}, nil
}

func (s *CommentService) Delete(ctx context.Context, req dto.DeleteCommentRequest) (*dto.CommentBoolResponse, error) {
	resp, err := s.rpc.CommentClient.DeleteComment(ctx, &commentpb.DeleteCommentReq{Id: req.ID, UserId: req.UserID})
	if err != nil {
		return nil, err
	}
	return &dto.CommentBoolResponse{Success: resp.GetSuccess()}, nil
}

func (s *CommentService) List(ctx context.Context, req dto.GetArticleCommentsRequest) (*dto.CommentListResponse, error) {
	resp, err := s.rpc.CommentClient.GetArticleComments(ctx, &commentpb.GetArticleCommentsReq{ArticleId: req.ArticleID, Page: req.Page, Size: req.Size})
	if err != nil {
		return nil, err
	}
	if err := s.attachLikeStatuses(ctx, req.UserID, resp.GetComments()); err != nil {
		return nil, err
	}
	return &dto.CommentListResponse{Comments: dto.ToCommentList(resp.GetComments()), Page: resp.GetPage(), Size: resp.GetSize()}, nil
}

func (s *CommentService) Replies(ctx context.Context, req dto.GetCommentRepliesRequest) (*dto.CommentRepliesResponse, error) {
	resp, err := s.rpc.CommentClient.GetCommentReplies(ctx, &commentpb.GetCommentRepliesReq{ParentId: req.ParentID, Page: req.Page, Size: req.Size})
	if err != nil {
		return nil, err
	}
	if err := s.attachLikeStatuses(ctx, req.UserID, resp.GetReplies()); err != nil {
		return nil, err
	}
	return &dto.CommentRepliesResponse{Replies: dto.ToCommentList(resp.GetReplies()), Page: resp.GetPage(), Size: resp.GetSize()}, nil
}

func (s *CommentService) attachLikeStatuses(ctx context.Context, userID uint64, comments []*commentpb.CommentInfo) error {
	objectIDs := make([]uint64, 0, len(comments))
	for _, comment := range comments {
		if comment != nil {
			objectIDs = append(objectIDs, comment.GetId())
		}
	}
	if len(objectIDs) == 0 {
		return nil
	}
	statusResp, err := s.rpc.LikeClient.BatchLikeStatus(ctx, &likepb.BatchCommentLikeStatusRequest{UserID: userID, ObjectType: "comment", ObjectIDs: objectIDs})
	if err != nil {
		return err
	}
	statuses := make(map[uint64]bool, len(statusResp.GetItems()))
	for _, item := range statusResp.GetItems() {
		if item != nil {
			statuses[item.GetObjectID()] = item.GetIsLiked()
		}
	}
	for _, comment := range comments {
		if comment != nil {
			comment.IsLiked = statuses[comment.GetId()]
		}
	}
	return nil
}
