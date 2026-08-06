package application

import (
	"context"

	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/likepb"
	"gateway/internal/model/dto"
)

type LikeService struct{ rpc *rpc.Client }

func NewLikeService(rpcClient *rpc.Client) *LikeService { return &LikeService{rpc: rpcClient} }

func (s *LikeService) ThumbUp(ctx context.Context, req dto.LikeRequest) (*dto.LikeBoolResponse, error) {
	resp, err := s.rpc.LikeClient.ThumbUp(ctx, &likepb.ThumbUpRequest{UserID: req.UserID, ObjectType: req.ObjectType, ObjectID: req.ObjectID})
	if err != nil {
		return nil, err
	}
	return &dto.LikeBoolResponse{Success: resp.GetSuccess()}, nil
}

func (s *LikeService) CancelThumbUp(ctx context.Context, req dto.LikeRequest) (*dto.LikeBoolResponse, error) {
	resp, err := s.rpc.LikeClient.CancelThumbUp(ctx, &likepb.CancelThumbUpRequest{UserID: req.UserID, ObjectType: req.ObjectType, ObjectID: req.ObjectID})
	if err != nil {
		return nil, err
	}
	return &dto.LikeBoolResponse{Success: resp.GetSuccess()}, nil
}

func (s *LikeService) HasLike(ctx context.Context, req dto.LikeStatusRequest) (*dto.LikeStatus, error) {
	resp, err := s.rpc.LikeClient.HasLike(ctx, &likepb.HasArticleLikeRequest{UserID: req.UserID, ObjectType: req.ObjectType, ObjectID: req.ObjectID})
	if err != nil {
		return nil, err
	}
	return &dto.LikeStatus{ObjectID: req.ObjectID, IsLiked: resp.GetIsLiked()}, nil
}

func (s *LikeService) BatchLikeStatus(ctx context.Context, req dto.BatchLikeStatusRequest) (*dto.BatchLikeStatusResponse, error) {
	resp, err := s.rpc.LikeClient.BatchLikeStatus(ctx, &likepb.BatchCommentLikeStatusRequest{UserID: req.UserID, ObjectType: req.ObjectType, ObjectIDs: req.ObjectIDs})
	if err != nil {
		return nil, err
	}
	items := make([]*dto.LikeStatus, 0, len(resp.GetItems()))
	for _, item := range resp.GetItems() {
		if item != nil {
			items = append(items, &dto.LikeStatus{ObjectID: item.GetObjectID(), IsLiked: item.GetIsLiked()})
		}
	}
	return &dto.BatchLikeStatusResponse{Items: items}, nil
}

func (s *LikeService) UserLikeList(ctx context.Context, req dto.UserLikeListRequest) (*likepb.PageQueryUserLikeListResponse, error) {
	return s.rpc.LikeClient.PageQueryUserLikeList(ctx, &likepb.PageQueryUserLikeListRequest{UserID: req.UserID, ObjectType: req.ObjectType, Page: req.Page, PageSize: req.PageSize})
}
