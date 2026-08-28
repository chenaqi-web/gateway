package application

import (
	"context"

	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/likepb"
	"gateway/internal/infras/clog"
	"gateway/internal/model/dto"

	"go.uber.org/zap"
)

type LikeService struct {
	rpc *rpc.Client
	log *clog.Log
}

func NewLikeService(rpcClient *rpc.Client, log *clog.Log) *LikeService {
	return &LikeService{rpc: rpcClient, log: log}
}

func (s *LikeService) ThumbUp(ctx context.Context, req dto.LikeRequest) (*dto.LikeBoolResponse, error) {
	resp, err := s.rpc.LikeClient.ThumbUp(ctx, &likepb.ThumbUpRequest{UserID: req.UserID, ObjectType: req.ObjectType, ObjectID: req.ObjectID})
	if err != nil {
		s.log.Error("LikeService/ThumbUp error", zap.Error(err))
		return nil, err
	}
	return dto.ToLikeBoolResponse(resp.GetSuccess()), nil
}

func (s *LikeService) CancelThumbUp(ctx context.Context, req dto.LikeRequest) (*dto.LikeBoolResponse, error) {
	resp, err := s.rpc.LikeClient.CancelThumbUp(ctx, &likepb.CancelThumbUpRequest{UserID: req.UserID, ObjectType: req.ObjectType, ObjectID: req.ObjectID})
	if err != nil {
		s.log.Error("LikeService/CancelThumbUp error", zap.Error(err))
		return nil, err
	}
	return dto.ToLikeBoolResponse(resp.GetSuccess()), nil
}

func (s *LikeService) UserLikeList(ctx context.Context, req dto.UserLikeListRequest) (*dto.UserLikeListResponse, error) {
	resp, err := s.rpc.LikeClient.PageQueryUserLikeList(ctx, &likepb.PageQueryUserLikeListRequest{UserID: req.UserID, ObjectType: req.ObjectType, Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		s.log.Error("LikeService/UserLikeList error", zap.Error(err))
		return nil, err
	}
	return dto.ToUserLikeListResponse(resp), nil
}
