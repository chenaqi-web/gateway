package dto

import "gateway/internal/client/rpc/core-rpc/likepb"

type LikeRequest struct {
	UserID     uint64 `json:"-"`
	ObjectType string `json:"objectType" binding:"required"`
	ObjectID   uint64 `json:"objectId" binding:"required"`
}

type LikeBoolResponse struct {
	Success bool `json:"success"`
}

// =====================================================================================================================

type UserLikeListRequest struct {
	UserID     uint64 `json:"-"`
	ObjectType string `json:"objectType" binding:"required"`
	Page       int32  `json:"page"`
	PageSize   int32  `json:"pageSize"`
}

type UserLikeListResponse struct {
	Articles []*Article `json:"articles"`
	Total    int64      `json:"total"`
}

func ToUserLikeListResponse(resp *likepb.PageQueryUserLikeListResponse) *UserLikeListResponse {
	if resp == nil {
		return &UserLikeListResponse{}
	}
	return &UserLikeListResponse{Articles: ToArticles(resp.GetArticles()), Total: resp.GetTotal()}
}
