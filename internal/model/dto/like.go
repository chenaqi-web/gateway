package dto

type LikeRequest struct {
	UserID     uint64 `json:"userId"`
	ObjectType string `json:"objectType" binding:"required"`
	ObjectID   uint64 `json:"objectId" binding:"required"`
}

type LikeStatusRequest struct {
	UserID     uint64 `json:"userId"`
	ObjectType string `json:"objectType" binding:"required"`
	ObjectID   uint64 `json:"objectId" binding:"required"`
}

type BatchLikeStatusRequest struct {
	UserID     uint64   `json:"userId"`
	ObjectType string   `json:"objectType" binding:"required"`
	ObjectIDs  []uint64 `json:"objectIds" binding:"required,min=1"`
}

type UserLikeListRequest struct {
	UserID     uint64 `json:"userId"`
	ObjectType string `json:"objectType" binding:"required"`
	Page       int32  `json:"page"`
	PageSize   int32  `json:"pageSize"`
}

type LikeBoolResponse struct {
	Success bool `json:"success"`
}

type LikeStatus struct {
	ObjectID uint64 `json:"objectId"`
	IsLiked  bool   `json:"isLiked"`
}

type BatchLikeStatusResponse struct {
	Items []*LikeStatus `json:"items"`
}
