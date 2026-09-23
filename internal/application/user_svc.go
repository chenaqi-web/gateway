package application

import (
	"context"
	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/userpb"
	"gateway/internal/config"
	"gateway/internal/infras/cache"
	"gateway/internal/infras/clog"
	"gateway/internal/infras/storage"
	"gateway/internal/model/dto"
	"gateway/internal/model/entity"
	"gateway/internal/utils"
	"time"

	"go.uber.org/zap"
)

type UserService struct {
	cfg           *config.Config
	rpc           *rpc.Client
	storage       *storage.Client
	userBlacklist *cache.Blacklist
	log           *clog.Log
}

func NewUserService(
	cfg *config.Config,
	rpcClient *rpc.Client,
	log *clog.Log,
	userBlacklist *cache.Blacklist,
	storage *storage.Client,
) *UserService {
	return &UserService{
		cfg:           cfg,
		rpc:           rpcClient,
		log:           log,
		userBlacklist: userBlacklist,
		storage:       storage,
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

func (s *UserService) UpdateProfile(ctx context.Context, req dto.UpdateProfileRequest) error {
	_, err := s.rpc.GetUserClient().UpdateProfile(ctx, &userpb.UpdateProfileRequest{
		UserId:   req.UserID,
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

func (s *UserService) UpdateAvatar(ctx context.Context, req *dto.UserAvatarRequest) (*dto.UserAvatarResponse, error) {
	// 1. 获取旧url
	old, err := s.rpc.GetUserClient().GetProfile(ctx, &userpb.GetProfileRequest{UserId: req.UserID})
	if err != nil {
		return nil, err
	}

	// 2.上传到文件到存储服务
	if err := utils.ValidateImage(s.cfg, req.File); err != nil {
		s.log.Error("StorageService/uploadImage error", zap.Error(err))
		return nil, err
	}

	result, err := s.storage.UploadAvatar(ctx, req.File, req.UserID)
	if err != nil {
		s.log.Error("StorageService/uploadImage error", zap.Error(err))
		return nil, err
	}

	// 3. 更新url
	resp, err := s.rpc.GetUserClient().UpdateAvatar(ctx, &userpb.UpdateAvatarRequest{
		UserId: req.UserID,
		Avatar: result,
	})
	if err != nil {
		_ = s.storage.Delete(ctx, result)
		s.log.Error("UserService/UpdateAvatar error", zap.Error(err))
		return nil, err
	}

	// 4.异步删除旧头像
	if old.GetAvatar() != "" && old.GetAvatar() != result {
		go func(key string) {
			deleteCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := s.storage.Delete(deleteCtx, key); err != nil {
				s.log.Warn("UserService/deleteOldAvatar error", zap.Error(err))
			}
		}(old.GetAvatar())
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

func (s *UserService) AddBlacklist(ctx context.Context, req *dto.UserBlacklistRequest) error {
	_, err := s.rpc.GetUserClient().UpdateUserStatus(ctx, &userpb.UpdateUserStatusRequest{
		UserId: req.UserID,
		Status: entity.StatusBlocked,
	})
	if err != nil {
		s.log.Error("UserService/AddBlacklist error", zap.Error(err))
		return err
	}

	err = s.userBlacklist.AddUser(ctx, req.UserID)
	if err != nil {
		s.log.Error("UserService/AddBlacklist error", zap.Error(err))
		return err
	}
	return nil
}

func (s *UserService) RemoveBlacklist(ctx context.Context, req *dto.UserBlacklistRequest) error {
	_, err := s.rpc.GetUserClient().UpdateUserStatus(ctx, &userpb.UpdateUserStatusRequest{
		UserId: req.UserID,
		Status: entity.StatusApproved,
	})
	if err != nil {
		s.log.Error("UserService/RemoveBlacklist error", zap.Error(err))
		return err
	}

	err = s.userBlacklist.RemoveUser(ctx, req.UserID)
	if err != nil {
		s.log.Error("UserService/RemoveBlacklist error", zap.Error(err))
		return err
	}
	return nil
}
