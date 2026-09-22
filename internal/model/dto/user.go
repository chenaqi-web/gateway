package dto

import (
	"gateway/internal/client/rpc/core-rpc/userpb"
)

type GetProfileResponse struct {
	ID                uint64 `json:"-"`
	Username          string `json:"username"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	Avatar            string `json:"avatar"`
	Sex               string `json:"sex"`
	Signature         string `json:"signature"`
	Birthday          string `json:"birthday"`
	Role              string `json:"role"`
	Status            string `json:"status"`
	ArticleCount      uint64 `json:"article_count"`
	FollowersCount    uint64 `json:"followers_count"`
	FollowingCount    uint64 `json:"following_count"`
	LikeCount         uint64 `json:"like_count"`
	ReceiveLikeCount  uint64 `json:"receive_like_count"`
	FavorCount        uint64 `json:"favor_count"`
	ReceiveFavorCount uint64 `json:"receive_favor_count"`
}
type UpdateProfileRequest struct {
	Username  string `json:"username" binding:"required,min=2,max=50"`
	Phone     string `json:"phone" binding:"max=20"`
	Sex       string `json:"sex" binding:"omitempty,oneof=male female"`
	Birthday  string `json:"birthday" binding:"omitempty,datetime=2006-01-02"`
	Signature string `json:"signature" binding:"max=255"`
}

type UpdateUserBlacklistRequest struct {
	UserID      uint64 `json:"user_id" binding:"required,gt=0"`
	Blacklisted bool   `json:"blacklisted"`
}

type UserAvatarRequest struct {
	UserID uint64 `json:"-"`
	Avatar string `json:"avatar"`
}
type UserAvatarResponse struct {
	Avatar string `json:"avatar"`
}

type UserInfo struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Sex      string `json:"sex"`
	Birthday string `json:"birthday"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

type UserListRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

type UserListResponse struct {
	Users []*UserInfo `json:"users"`
	Total uint64      `json:"total"`
}

type UserSearchRequest struct {
	Keyword  string `json:"keyword"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

type UserSearchResponse struct {
	Users []*UserInfo `json:"users"`
	Total uint64      `json:"total"`
}

// =====================================================================================================================

func ToUserInfo(user *userpb.UserInfo) *UserInfo {
	if user == nil {
		return nil
	}
	return &UserInfo{
		ID:       user.GetId(),
		Username: user.GetUsername(),
		Email:    user.GetEmail(),
		Phone:    user.GetPhone(),
		Avatar:   user.GetAvatar(),
		Sex:      user.GetSex(),
		Role:     user.GetRole(),
		Status:   user.GetStatus(),
		Birthday: user.GetBirthday(),
	}
}

func ToUserListResponse(users []*userpb.UserInfo, total uint64) *UserListResponse {
	if users == nil {
		return nil
	}
	res := make([]*UserInfo, 0, len(users))
	for _, user := range users {
		res = append(res, ToUserInfo(user))
	}
	return &UserListResponse{
		Users: res,
		Total: total,
	}
}

func ToUserSearchResponse(users []*userpb.UserInfo, total uint64) *UserSearchResponse {
	if users == nil {
		return nil
	}
	res := make([]*UserInfo, 0, len(users))
	for _, user := range users {
		res = append(res, ToUserInfo(user))
	}
	return &UserSearchResponse{
		Users: res,
		Total: total,
	}
}

func ToUserAvatarResponse(url string) *UserAvatarResponse {
	return &UserAvatarResponse{
		Avatar: url,
	}
}

func ToGetProfileResponse(user *userpb.GetProfileResponse) *GetProfileResponse {
	return &GetProfileResponse{
		ID:                user.Id,
		Username:          user.Username,
		Email:             user.Email,
		Phone:             user.Phone,
		Avatar:            user.Avatar,
		Sex:               user.Sex,
		Signature:         user.Signature,
		Birthday:          user.Birthday,
		Role:              user.Role,
		Status:            user.Status,
		ArticleCount:      user.ArticleCount,
		FollowersCount:    user.FollowersCount,
		FollowingCount:    user.FollowingCount,
		LikeCount:         user.LikeCount,
		ReceiveLikeCount:  user.ReceiveLikeCount,
		FavorCount:        user.FavorCount,
		ReceiveFavorCount: user.ReceiveFavorCount,
	}
}
