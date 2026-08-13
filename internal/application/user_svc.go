package application

import (
	"context"
	"errors"
	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/userpb"
	"gateway/internal/infras/cache"
	"gateway/internal/model/dto"
)

type UserService struct {
	rpc          *rpc.Client
	jwtBlacklist *cache.JwtBlacklist
}

func NewUserService(rpcClient *rpc.Client, jwtBlacklist *cache.JwtBlacklist) *UserService {
	return &UserService{
		rpc:          rpcClient,
		jwtBlacklist: jwtBlacklist,
	}
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

func (s *UserService) List(ctx context.Context, keyword string, page, pageSize uint32) ([]*dto.UserProfile, uint64, error) {
	resp, err := s.rpc.GetUserClient().ListUsers(ctx, &userpb.ListUsersRequest{Keyword: keyword, Page: page, PageSize: pageSize})
	if err != nil {
		return nil, 0, err
	}
	users := make([]*dto.UserProfile, 0, len(resp.GetUsers()))
	for _, user := range resp.GetUsers() {
		users = append(users, dto.ToUserProfile(user))
	}
	return users, resp.GetTotal(), nil
}

func (s *UserService) UpdateStatus(ctx context.Context, userID uint64, userStatus string) (bool, error) {
	// 1.加入黑名单，首先将token存入redis

	// 2. 修改数据库状态
	resp, err := s.rpc.GetUserClient().UpdateUserStatus(ctx, &userpb.UpdateUserStatusRequest{UserId: userID, Status: userStatus})
	if err != nil {
		return false, err
	}
	if !resp.GetSuccess() {
		return false, errors.New("user status update failed")
	}
	return true, nil
}
