package application

import (
	"context"

	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/articlepb"
	"gateway/internal/config"
	"gateway/internal/infras/clog"
	"gateway/internal/infras/storage"
	"gateway/internal/model/dto"
	"gateway/internal/utils"

	"go.uber.org/zap"
)

type ArticleService struct {
	cfg     *config.Config
	rpc     *rpc.Client
	storage *storage.Client
	log     *clog.Log
}

func NewArticleService(cfg *config.Config, rpcClient *rpc.Client, log *clog.Log, storageClient *storage.Client) *ArticleService {
	return &ArticleService{
		cfg:     cfg,
		rpc:     rpcClient,
		storage: storageClient,
		log:     log,
	}
}

func (s *ArticleService) Create(ctx context.Context, req dto.CreateArticleRequest) (*dto.ArticleBoolResponse, error) {
	resp, err := s.rpc.ArticleClient.CreateArticle(ctx, &articlepb.CreateArticleRequest{
		AuthorID:    req.AuthorID,
		Title:       req.Title,
		Summary:     req.Summary,
		Content:     req.Content,
		CoverImage:  req.CoverImage,
		CategoryID:  req.CategoryID,
		IsTop:       req.IsTop,
		IsPublished: req.IsPublish,
	})
	if err != nil {
		s.log.Error("ArticleService/Create error", zap.Error(err))
		return nil, err
	}
	return dto.ToArticleBoolResponse(resp.GetSuccess()), nil
}

func (s *ArticleService) Edit(ctx context.Context, req dto.EditArticleRequest) (*dto.EditArticleResponse, error) {
	resp, err := s.rpc.ArticleClient.EditorArticle(ctx, &articlepb.EditorArticleRequest{
		Id:          req.ID,
		AuthorID:    req.AuthorID,
		Title:       req.Title,
		Summary:     req.Summary,
		Content:     req.Content,
		CoverImage:  req.CoverImage,
		CategoryID:  req.CategoryID,
		IsTop:       req.IsTop,
		IsPublished: req.IsPublish,
	})
	if err != nil {
		s.log.Error("ArticleService/Edit error", zap.Error(err))
		return nil, err
	}
	return dto.ToEditArticleResponse(resp.GetSuccess(), resp.GetArticleID()), nil
}

func (s *ArticleService) PublishDraft(ctx context.Context, req dto.PublishDraftRequest) (*dto.ArticleBoolResponse, error) {
	resp, err := s.rpc.ArticleClient.PublishDraft(ctx, &articlepb.PublishDraftRequest{
		Id:       req.ID,
		AuthorID: req.AuthorID,
	})
	if err != nil {
		s.log.Error("ArticleService/PublishDraft error", zap.Error(err))
		return nil, err
	}
	return dto.ToArticleBoolResponse(resp.GetSuccess()), nil
}

func (s *ArticleService) Search(ctx context.Context, req dto.SearchArticlesRequest) (*dto.ListArticlesResponse, error) {
	resp, err := s.rpc.ArticleClient.SearchArticles(ctx, &articlepb.SearchArticlesRequest{Q: req.Q, Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		s.log.Error("ArticleService/Search error", zap.Error(err))
		return nil, err
	}
	return dto.ToListArticlesResponse(resp.GetArticles()), nil
}

func (s *ArticleService) Delete(ctx context.Context, req dto.DeleteArticleRequest) (*dto.ArticleBoolResponse, error) {
	resp, err := s.rpc.ArticleClient.DeleteArticle(ctx, &articlepb.DeleteArticleRequest{Id: req.ID, UserID: req.AuthorID, Role: req.Role})
	if err != nil {
		s.log.Error("ArticleService/Delete error", zap.Error(err))
		return nil, err
	}
	return dto.ToArticleBoolResponse(resp.GetSuccess()), nil
}

func (s *ArticleService) DeleteDraft(ctx context.Context, req dto.DeleteDraftRequest) (*dto.ArticleBoolResponse, error) {
	resp, err := s.rpc.ArticleClient.DeleteDraft(ctx, &articlepb.DeleteDraftRequest{Id: req.ID, AuthorID: req.AuthorID})
	if err != nil {
		s.log.Error("ArticleService/DeleteDraft error", zap.Error(err))
		return nil, err
	}
	return dto.ToArticleBoolResponse(resp.GetSuccess()), nil
}

// =====================================================================================================================

func (s *ArticleService) UploadCover(ctx context.Context, req dto.ArticleImageUploadRequest) (*dto.ArticleImageUploadResponse, error) {
	if err := utils.ValidateImage(s.cfg, req.File); err != nil {
		s.log.Error("ArticleService/UploadCover validate error", zap.Error(err))
		return nil, err
	}

	url, err := s.storage.Upload(ctx, req.File, storage.DirectoryArticleCover, req.UserID)
	if err != nil {
		s.log.Error("ArticleService/UploadCover upload error", zap.Error(err))
		return nil, err
	}
	return dto.ToArticleImageUploadResponse(url), nil
}

func (s *ArticleService) UploadContentImage(ctx context.Context, req dto.ArticleImageUploadRequest) (*dto.ArticleImageUploadResponse, error) {
	if err := utils.ValidateImage(s.cfg, req.File); err != nil {
		s.log.Error("ArticleService/UploadContentImage validate error", zap.Error(err))
		return nil, err
	}

	url, err := s.storage.Upload(ctx, req.File, storage.DirectoryArticleContent, req.UserID)
	if err != nil {
		s.log.Error("ArticleService/UploadContentImage upload error", zap.Error(err))
		return nil, err
	}
	return dto.ToArticleImageUploadResponse(url), nil
}

// =====================================================================================================================

func (s *ArticleService) GetDetail(ctx context.Context, req dto.GetArticleRequest) (*dto.GetArticleResponse, error) {
	resp, err := s.rpc.ArticleClient.GetArticle(ctx, &articlepb.GetArticleRequest{Id: req.ID})
	if err != nil {
		s.log.Error("ArticleService/GetDetail error", zap.Error(err))
		return nil, err
	}
	return dto.ToGetArticleResponse(resp), nil
}

func (s *ArticleService) List(ctx context.Context, req dto.ListArticlesRequest) (*dto.ListArticlesResponse, error) {
	resp, err := s.rpc.ArticleClient.ListArticles(ctx, &articlepb.ListArticlesRequest{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		s.log.Error("ArticleService/List error", zap.Error(err))
		return nil, err
	}
	return dto.ToListArticlesResponse(resp.GetArticles()), nil
}

func (s *ArticleService) ListByUserID(ctx context.Context, req dto.ListMyArticlesRequest) (*dto.ListArticlesResponse, error) {
	resp, err := s.rpc.ArticleClient.ListMyArticles(ctx, &articlepb.ListMyArticlesRequest{
		AuthorID:    req.AuthorID,
		Page:        req.Page,
		PageSize:    req.PageSize,
		IsPublished: req.IsPublished,
	})
	if err != nil {
		s.log.Error("ArticleService/ListByUserID error", zap.Error(err))
		return nil, err
	}
	return dto.ToListMyArticlesResponse(resp), nil
}

func (s *ArticleService) ListByCategory(ctx context.Context, req dto.ListByCategoryRequest) (*dto.ListArticlesResponse, error) {
	resp, err := s.rpc.ArticleClient.ListByCategory(ctx, &articlepb.ListByCategoryRequest{CategoryID: req.CategoryID, Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		s.log.Error("ArticleService/ListByCategory error", zap.Error(err))
		return nil, err
	}
	return dto.ToListArticlesResponse(resp.GetArticles()), nil
}
