package application

import (
	"context"
	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/userpb"
	"gateway/internal/infras/cache"
	"gateway/internal/infras/clog"
	"gateway/internal/model/dto"
	"gateway/internal/model/entity"

	"go.uber.org/zap"
)

type UserService struct {
	rpc           *rpc.Client
	log           *clog.Log
	userBlacklist *cache.Blacklist
}

func NewUserService(
	rpcClient *rpc.Client,
	log *clog.Log,
	userBlacklist *cache.Blacklist) *UserService {
	return &UserService{
		rpc:           rpcClient,
		log:           log,
		userBlacklist: userBlacklist,
	}
}

func (s *UserService) GetProfile(ctx context.Context, userID uint64) (*dto.UserProfile, error) {
	resp, err := s.rpc.GetUserClient().GetProfile(ctx, &userpb.GetProfileRequest{UserId: userID})
	if err != nil {
		s.log.Error("UserService/GetProfile error", zap.Error(err))
		return nil, err
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
		s.log.Error("UserService/UpdateProfile error", zap.Error(err))
		return nil, err
	}
	return dto.ToUserProfile(resp.GetUser()), nil
}

func (s *UserService) UpdateAvatar(ctx context.Context, userID uint64, avatar string) (*dto.UserProfile, error) {
	resp, err := s.rpc.GetUserClient().UpdateAvatar(ctx, &userpb.UpdateAvatarRequest{
		UserId: userID,
		Avatar: avatar,
	})
	if err != nil {
		s.log.Error("UserService/UpdateAvatar error", zap.Error(err))
		return nil, err
	}
	return dto.ToUserProfile(resp.GetUser()), nil
}

func (s *UserService) UserList(ctx context.Context, keyword string, page, pageSize uint32) ([]*dto.UserProfile, uint64, error) {
	resp, err := s.rpc.GetUserClient().ListUsers(ctx, &userpb.ListUsersRequest{
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		s.log.Error("UserService/List error", zap.Error(err))
		return nil, 0, err
	}
	users := make([]*dto.UserProfile, 0, len(resp.GetUsers()))
	for _, user := range resp.GetUsers() {
		users = append(users, dto.ToUserProfile(user))
	}
	return users, resp.GetTotal(), nil
}

func (s *UserService) UpdateBlacklist(ctx context.Context, userID uint64, blacklisted bool) (*dto.UserBoolResponse, error) {
	userStatus := entity.StatusApproved
	if blacklisted {
		userStatus = entity.StatusBlocked
	}

	resp, err := s.rpc.GetUserClient().UpdateUserStatus(ctx, &userpb.UpdateUserStatusRequest{UserId: userID, Status: userStatus})
	if err != nil {
		s.log.Error("UserService/UpdateBlacklist error", zap.Error(err))
		return nil, err
	}

	if blacklisted {
		err = s.userBlacklist.AddUser(ctx, userID)
	} else {
		err = s.userBlacklist.RemoveUser(ctx, userID)
	}
	if err != nil {
		s.log.Error("UserService/UpdateBlacklist error", zap.Error(err))
		return nil, err
	}
	return dto.ToUserBoolResponse(resp.GetSuccess()), nil
}
