package dto

import "gateway/internal/client/rpc/core-rpc/articlepb"

// ---------- 实体 ----------

type Article struct {
	ID           uint64 `json:"id"`
	Title        string `json:"title"`
	Summary      string `json:"summary"`
	Content      string `json:"content"`
	CoverImage   string `json:"coverImage"`
	AuthorID     uint64 `json:"authorID"`
	CategoryID   uint64 `json:"categoryID"`
	IsTop        bool   `json:"isTop"`
	ViewCount    uint64 `json:"viewCount"`
	LikeCount    uint64 `json:"likeCount"`
	CommentCount uint64 `json:"commentCount"`
	CreatedAt    uint64 `json:"createdAt"`
	UpdatedAt    uint64 `json:"updatedAt"`
	AuthorName   string `json:"authorName"`
	AuthorAvatar string `json:"authorAvatar"`
}

// ---------- 请求 ----------

type CreateArticleRequest struct {
	AuthorID   uint64 `json:"-"`
	CategoryID uint64 `json:"categoryID" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Summary    string `json:"summary"`
	CoverImage string `json:"coverImage"`
	IsTop      bool   `json:"isTop"`
}

type GetArticleRequest struct {
	ID uint64 `json:"id" binding:"required"`
}

type ListArticlesRequest struct {
	Page     uint32 `json:"page"`
	PageSize uint32 `json:"pageSize"`
}

type ListMyArticlesRequest struct {
	AuthorID uint64 `json:"-"`
	Page     uint32 `json:"page"`
	PageSize uint32 `json:"pageSize"`
}

type ListByCategoryRequest struct {
	CategoryID uint64 `json:"categoryID" binding:"required"`
	Page       uint32 `json:"page"`
	PageSize   uint32 `json:"pageSize"`
}

type SearchArticlesRequest struct {
	Q        string `json:"q" binding:"required"`
	Page     uint32 `json:"page"`
	PageSize uint32 `json:"pageSize"`
}

type DeleteArticleRequest struct {
	ID       uint64 `json:"id" binding:"required"`
	AuthorID uint64 `json:"-"`
}

// ---------- 响应 ----------

type ArticleBoolResponse struct {
	Success bool `json:"success"`
}

type GetArticleResponse struct {
	Article *Article `json:"article"`
}

type ListArticlesResponse struct {
	Articles []*Article `json:"articles"`
}

// ---------- 转换 ----------

func ToArticle(item *articlepb.Article) *Article {
	if item == nil {
		return nil
	}
	return &Article{
		ID:           item.GetId(),
		Title:        item.GetTitle(),
		Summary:      item.GetSummary(),
		Content:      item.GetContent(),
		CoverImage:   item.GetCoverImage(),
		AuthorID:     item.GetAuthorID(),
		CategoryID:   item.GetCategoryID(),
		IsTop:        item.GetIsTop(),
		ViewCount:    item.GetViewCount(),
		LikeCount:    item.GetLikeCount(),
		CommentCount: item.GetCommentCount(),
		CreatedAt:    item.GetCreatedAt(),
		UpdatedAt:    item.GetUpdatedAt(),
		AuthorName:   item.GetAuthorName(),
		AuthorAvatar: item.GetAuthorAvatar(),
	}
}

func ToArticles(items []*articlepb.Article) []*Article {
	list := make([]*Article, 0, len(items))
	for _, item := range items {
		list = append(list, ToArticle(item))
	}
	return list
}

func ToArticleBoolResponse(success bool) *ArticleBoolResponse {
	return &ArticleBoolResponse{Success: success}
}

func ToGetArticleResponse(resp *articlepb.GetArticleResponse) *GetArticleResponse {
	if resp == nil {
		return &GetArticleResponse{}
	}
	return &GetArticleResponse{Article: ToArticle(resp.GetArticle())}
}

func ToListArticlesResponse(articles []*articlepb.Article) *ListArticlesResponse {
	return &ListArticlesResponse{Articles: ToArticles(articles)}
}
