package dto

import (
	"gateway/internal/client/rpc/core-rpc/articlepb"
	"mime/multipart"
)

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
	IsPublished  bool   `json:"isPublished"`
	ViewCount    uint64 `json:"viewCount"`
	LikeCount    uint64 `json:"likeCount"`
	FavorCount   uint64 `json:"favorCount"`
	CommentCount uint64 `json:"commentCount"`
	CreatedAt    uint64 `json:"createdAt"`
	UpdatedAt    uint64 `json:"updatedAt"`
	PublishedAt  uint64 `json:"publishedAt"`
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
	IsPublish  bool   `json:"isPublish"`
}

type EditArticleRequest struct {
	ID         uint64 `json:"id" binding:"required"`
	AuthorID   uint64 `json:"-"`
	CategoryID uint64 `json:"categoryID" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Summary    string `json:"summary"`
	CoverImage string `json:"coverImage"`
	IsTop      bool   `json:"isTop"`
	IsPublish  bool   `json:"isPublish"`
}

type GetArticleRequest struct {
	ID uint64 `json:"id" binding:"required"`
}

type ListArticlesRequest struct {
	Page     uint32 `json:"page"`
	PageSize uint32 `json:"pageSize"`
}

type ListMyArticlesRequest struct {
	AuthorID    uint64 `json:"authorID"`
	Page        uint32 `json:"page"`
	PageSize    uint32 `json:"pageSize"`
	IsPublished *bool  `json:"isPublished"`
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
	Role     string `json:"-"`
}

type PublishDraftRequest struct {
	ID       uint64 `json:"id" binding:"required"`
	AuthorID uint64 `json:"-"`
}

type DeleteDraftRequest struct {
	ID       uint64 `json:"id" binding:"required"`
	AuthorID uint64 `json:"-"`
}

type ArticleImageUploadRequest struct {
	UserID uint64                `json:"-"`
	File   *multipart.FileHeader `json:"-"`
}

// ---------- 响应 ----------

type ArticleBoolResponse struct {
	Success bool `json:"success"`
}

type EditArticleResponse struct {
	Success   bool   `json:"success"`
	ArticleID uint64 `json:"articleID"`
}

type ArticleImageUploadResponse struct {
	URL string `json:"url"`
}

type GetArticleResponse struct {
	Article *Article `json:"article"`
	IsLiked bool     `json:"isLiked"`
}

type ListArticlesResponse struct {
	Articles []*Article `json:"articles"`
	Total    uint64     `json:"total"`
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
		IsPublished:  item.GetIsPublished(),
		ViewCount:    item.GetViewCount(),
		LikeCount:    item.GetLikeCount(),
		FavorCount:   item.GetFavorCount(),
		CommentCount: item.GetCommentCount(),
		CreatedAt:    item.GetCreatedAt(),
		UpdatedAt:    item.GetUpdatedAt(),
		PublishedAt:  item.GetPublishedAt(),
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

func ToEditArticleResponse(success bool, articleID uint64) *EditArticleResponse {
	return &EditArticleResponse{Success: success, ArticleID: articleID}
}

func ToArticleImageUploadResponse(url string) *ArticleImageUploadResponse {
	return &ArticleImageUploadResponse{URL: url}
}

func ToGetArticleResponse(resp *articlepb.GetArticleResponse, isLiked bool) *GetArticleResponse {
	if resp == nil {
		return &GetArticleResponse{}
	}
	return &GetArticleResponse{
		Article: &Article{
			ID:           resp.GetId(),
			Title:        resp.GetTitle(),
			Summary:      resp.GetSummary(),
			Content:      resp.GetContent(),
			CoverImage:   resp.GetCoverImage(),
			AuthorID:     resp.GetAuthorID(),
			CategoryID:   resp.GetCategoryID(),
			IsTop:        resp.GetIsTop(),
			IsPublished:  resp.GetIsPublished(),
			ViewCount:    resp.GetViewCount(),
			LikeCount:    resp.GetLikeCount(),
			FavorCount:   resp.GetFavorCount(),
			CommentCount: resp.GetCommentCount(),
			CreatedAt:    resp.GetCreatedAt(),
			UpdatedAt:    resp.GetUpdatedAt(),
			PublishedAt:  resp.GetPublishedAt(),
			AuthorName:   resp.GetAuthorName(),
			AuthorAvatar: resp.GetAuthorAvatar(),
		},
		IsLiked: isLiked,
	}
}

func ToListArticlesResponse(articles []*articlepb.Article) *ListArticlesResponse {
	return &ListArticlesResponse{Articles: ToArticles(articles), Total: uint64(len(articles))}
}

func ToListMyArticlesResponse(resp *articlepb.ListMyArticlesResponse) *ListArticlesResponse {
	if resp == nil {
		return &ListArticlesResponse{}
	}
	return &ListArticlesResponse{Articles: ToArticles(resp.GetArticles()), Total: resp.GetTotal()}
}
