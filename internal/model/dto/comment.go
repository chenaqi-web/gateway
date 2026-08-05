package dto

import "gateway/internal/client/rpc/core-rpc/commentpb"

type CreateCommentRequest struct {
	ArticleID uint64 `json:"articleId" binding:"required"`
	UserID    uint64 `json:"userId" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

type CreateReplyRequest struct {
	ArticleID uint64 `json:"articleId" binding:"required"`
	ParentID  uint64 `json:"parentId" binding:"required"`
	UserID    uint64 `json:"userId" binding:"required"`
	ReplyToID uint64 `json:"replyToId"`
	Content   string `json:"content" binding:"required"`
}

type DeleteCommentRequest struct {
	ID     uint64 `json:"id" binding:"required"`
	UserID uint64 `json:"userId" binding:"required"`
}

type GetArticleCommentsRequest struct {
	ArticleID uint64 `json:"articleId" binding:"required"`
	Page      int32  `json:"page"`
	Size      int32  `json:"size"`
}

type GetCommentRepliesRequest struct {
	ParentID uint64 `json:"parentId" binding:"required"`
	Page     int32  `json:"page"`
	Size     int32  `json:"size"`
}

type CommentInfo struct {
	ID         uint64 `json:"id"`
	ArticleID  uint64 `json:"articleId"`
	UserID     uint64 `json:"userId"`
	ParentID   uint64 `json:"parentId"`
	RootID     uint64 `json:"rootId"`
	ReplyToID  uint64 `json:"replyToId"`
	Content    string `json:"content"`
	LikeCount  uint32 `json:"likeCount"`
	ChildCount uint32 `json:"childCount"`
	CreatedAt  string `json:"createdAt"`
	UserName   string `json:"userName"`
	UserAvatar string `json:"userAvatar"`
	IsLiked    bool   `json:"isLiked"`
}

type CommentListResponse struct {
	Comments []*CommentInfo `json:"comments"`
	Page     int32          `json:"page"`
	Size     int32          `json:"size"`
}

type CommentRepliesResponse struct {
	Replies []*CommentInfo `json:"replies"`
	Page    int32          `json:"page"`
	Size    int32          `json:"size"`
}

type CommentBoolResponse struct {
	Success bool `json:"success"`
}

func ToCommentInfo(item *commentpb.CommentInfo) *CommentInfo {
	if item == nil {
		return nil
	}
	return &CommentInfo{ID: item.GetId(), ArticleID: item.GetArticleId(), UserID: item.GetUserId(), ParentID: item.GetParentId(), RootID: item.GetRootId(), ReplyToID: item.GetReplyToId(), Content: item.GetContent(), LikeCount: item.GetLikeCount(), ChildCount: item.GetChildCount(), CreatedAt: item.GetCreatedAt(), UserName: item.GetUserName(), UserAvatar: item.GetUserAvatar(), IsLiked: item.GetIsLiked()}
}

func ToCommentList(items []*commentpb.CommentInfo) []*CommentInfo {
	result := make([]*CommentInfo, 0, len(items))
	for _, item := range items {
		result = append(result, ToCommentInfo(item))
	}
	return result
}
