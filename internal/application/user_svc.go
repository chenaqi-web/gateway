package application

import (
	"context"
	"errors"

	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/userpb"
	"gateway/internal/model/dto"
)

type UserService struct{ rpc *rpc.Client }

func NewUserService(rpcClient *rpc.Client) *UserService { return &UserService{rpc: rpcClient} }

func (s *UserService) Get(ctx context.Context) (any, error) {
	return s.rpc.GetUserClient().Login(ctx, &userpb.LoginReq{})
}

func (s *UserService) GetProfile(ctx context.Context, userID uint64) (*dto.UserProfile, error) {
	resp, err := s.rpc.GetUserClient().GetProfile(ctx, &userpb.GetProfileRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	if resp.GetUser() == nil {
		return nil, errors.New("invalid core response")
	}
	return dto.ToUserProfile(resp.GetUser()), nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uint64, req dto.UpdateProfileRequest) (*dto.UserProfile, error) {
	resp, err := s.rpc.GetUserClient().UpdateProfile(ctx, &userpb.UpdateProfileRequest{
		UserId:   userID,
		Username: req.Username,
		Phone:    req.Phone,
		Sex:      req.Sex,
		Age:      req.Age,
	})
	if err != nil {
		return nil, err
	}
	if resp.GetUser() == nil {
		return nil, errors.New("invalid core response")
	}
	return dto.ToUserProfile(resp.GetUser()), nil
}

func (s *UserService) UpdateAvatar(ctx context.Context, userID uint64, avatar string) (*dto.UserProfile, error) {
	resp, err := s.rpc.GetUserClient().UpdateAvatar(ctx, &userpb.UpdateAvatarRequest{
		UserId: userID,
		Avatar: avatar,
	})
	if err != nil {
		return nil, err
	}
	if resp.GetUser() == nil {
		return nil, errors.New("invalid core response")
	}
	return dto.ToUserProfile(resp.GetUser()), nil
}
