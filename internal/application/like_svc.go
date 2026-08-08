package application

import (
	"context"

	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/likepb"
	"gateway/internal/model/dto"
)

type LikeService struct{ rpc *rpc.Client }

func NewLikeService(rpcClient *rpc.Client) *LikeService { return &LikeService{rpc: rpcClient} }

func (s *LikeService) ThumbUp(ctx context.Context, userID uint64, req dto.LikeRequest) (*dto.LikeBoolResponse, error) {
	resp, err := s.rpc.LikeClient.ThumbUp(ctx, &likepb.ThumbUpRequest{UserID: userID, ObjectType: req.ObjectType, ObjectID: req.ObjectID})
	if err != nil {
		return nil, err
	}
	return &dto.LikeBoolResponse{Success: resp.GetSuccess()}, nil
}

func (s *LikeService) CancelThumbUp(ctx context.Context, userID uint64, req dto.LikeRequest) (*dto.LikeBoolResponse, error) {
	resp, err := s.rpc.LikeClient.CancelThumbUp(ctx, &likepb.CancelThumbUpRequest{UserID: userID, ObjectType: req.ObjectType, ObjectID: req.ObjectID})
	if err != nil {
		return nil, err
	}
	return &dto.LikeBoolResponse{Success: resp.GetSuccess()}, nil
}

func (s *LikeService) UserLikeList(ctx context.Context, userID uint64, req dto.UserLikeListRequest) (*likepb.PageQueryUserLikeListResponse, error) {
	return s.rpc.LikeClient.PageQueryUserLikeList(ctx, &likepb.PageQueryUserLikeListRequest{UserID: userID, ObjectType: req.ObjectType, Page: req.Page, PageSize: req.PageSize})
}

func (s *LikeService) HasLike(ctx context.Context, userID uint64, req dto.LikeStatusRequest) (*dto.LikeStatus, error) {
	resp, err := s.rpc.LikeClient.HasLike(ctx, &likepb.HasArticleLikeRequest{UserID: userID, ObjectType: req.ObjectType, ObjectID: req.ObjectID})
	if err != nil {
		return nil, err
	}
	return &dto.LikeStatus{ObjectID: req.ObjectID, IsLiked: resp.GetIsLiked()}, nil
}
