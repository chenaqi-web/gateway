package dto

import "gateway/internal/client/rpc/core-rpc/userpb"

type UserProfile struct {
	ID               uint64 `json:"id"`
	Username         string `json:"username"`
	Email            string `json:"email"`
	Phone            string `json:"phone"`
	Avatar           string `json:"avatar"`
	Sex              string `json:"sex"`
	Age              uint32 `json:"age"`
	Role             string `json:"role"`
	Status           string `json:"status"`
	LikeCount        uint64 `json:"like_count"`
	ReceiveLikeCount uint64 `json:"receive_like_count"`
}

type UpdateProfileRequest struct {
	Username string `json:"username" binding:"required,min=2,max=50"`
	Phone    string `json:"phone" binding:"max=20"`
	Sex      string `json:"sex" binding:"omitempty,oneof=male female"`
	Age      uint32 `json:"age" binding:"max=150"`
}

type UpdateUserStatusRequest struct {
	UserID uint64 `json:"user_id" binding:"required,gt=0"`
	Status string `json:"status" binding:"required,oneof=approved blocked"`
}

func ToUserProfile(user *userpb.UserInfo) *UserProfile {
	if user == nil {
		return nil
	}
	return &UserProfile{
		ID:               user.GetId(),
		Username:         user.GetUsername(),
		Email:            user.GetEmail(),
		Phone:            user.GetPhone(),
		Avatar:           user.GetAvatar(),
		Sex:              user.GetSex(),
		Age:              user.GetAge(),
		Role:             user.GetRole(),
		Status:           user.GetStatus(),
		LikeCount:        user.GetLikeCount(),
		ReceiveLikeCount: user.GetReceiveLikeCount(),
	}
}
