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

func (s *UserService) GetProfile(ctx context.Context, userID uint64) (*dto.GetProfileResponse, error) {
	res, err := s.rpc.GetUserClient().GetProfile(ctx, &userpb.GetProfileRequest{UserId: userID})
	if err != nil {
		s.log.Error("UserService/GetProfile error", zap.Error(err))
		return nil, err
	}
	return dto.ToGetProfileResponse(res), nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uint64, req dto.UpdateProfileRequest) error {
	_, err := s.rpc.GetUserClient().UpdateProfile(ctx, &userpb.UpdateProfileRequest{
		UserId:   userID,
		Username: req.Username,
		Phone:    req.Phone,
		Sex:      req.Sex,
		Birthday: req.Birthday,
	})
	if err != nil {
		s.log.Error("UserService/UpdateProfile error", zap.Error(err))
		return err
	}
	return nil
}

// todo 图床的内容后续修改

func (s *UserService) UpdateAvatar(ctx context.Context, req dto.UserAvatarRequest) (*dto.UserAvatarResponse, error) {
	resp, err := s.rpc.GetUserClient().UpdateAvatar(ctx, &userpb.UpdateAvatarRequest{
		UserId: req.UserID,
		Avatar: req.Avatar,
	})
	if err != nil {
		s.log.Error("UserService/UpdateAvatar error", zap.Error(err))
		return nil, err
	}
	return dto.ToUserAvatarResponse(resp.Url), nil
}

// =====================================================================================================================

func (s *UserService) SearchUser(ctx context.Context, req *dto.UserSearchRequest) (*dto.UserSearchResponse, error) {
	resp, err := s.rpc.GetUserClient().SearchUsers(ctx, &userpb.SearchUsersRequest{
		Keyword:  req.Keyword,
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	})
	if err != nil {
		s.log.Error("UserService/SearchUser error", zap.Error(err))
		return nil, err
	}
	return dto.ToUserSearchResponse(resp.GetUsers(), resp.GetTotal()), nil
}

func (s *UserService) UserList(ctx context.Context, rep *dto.UserListRequest) (*dto.UserListResponse, error) {
	resp, err := s.rpc.GetUserClient().ListUsers(ctx, &userpb.ListUsersRequest{
		Page:     int32(rep.Page),
		PageSize: int32(rep.PageSize),
	})
	if err != nil {
		s.log.Error("UserService/List error", zap.Error(err))
		return nil, err
	}
	return dto.ToUserListResponse(resp.GetUsers(), resp.GetTotal()), nil
}

func (s *UserService) UpdateBlacklist(ctx context.Context, userID uint64, blacklisted bool) error {
	userStatus := entity.StatusApproved
	if blacklisted {
		userStatus = entity.StatusBlocked
	}

	_, err := s.rpc.GetUserClient().UpdateUserStatus(ctx, &userpb.UpdateUserStatusRequest{UserId: userID, Status: userStatus})
	if err != nil {
		s.log.Error("UserService/UpdateBlacklist error", zap.Error(err))
		return err
	}

	if blacklisted {
		err = s.userBlacklist.AddUser(ctx, userID)
	} else {
		err = s.userBlacklist.RemoveUser(ctx, userID)
	}
	if err != nil {
		s.log.Error("UserService/UpdateBlacklist error", zap.Error(err))
		return err
	}
	return nil
}
