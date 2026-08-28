package application

import (
	"context"

	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/articlepb"
	"gateway/internal/infras/clog"
	"gateway/internal/model/dto"

	"go.uber.org/zap"
)

type ArticleService struct {
	rpc *rpc.Client
	log *clog.Log
}

func NewArticleService(rpcClient *rpc.Client, log *clog.Log) *ArticleService {
	return &ArticleService{rpc: rpcClient, log: log}
}

func (s *ArticleService) Create(ctx context.Context, req dto.CreateArticleRequest) (*dto.ArticleBoolResponse, error) {
	resp, err := s.rpc.ArticleClient.CreateArticle(ctx, &articlepb.CreateArticleRequest{
		AuthorID: req.AuthorID, Title: req.Title, Summary: req.Summary, Content: req.Content,
		CoverImage: req.CoverImage, CategoryID: req.CategoryID, IsTop: req.IsTop,
	})
	if err != nil {
		s.log.Error("ArticleService/Create error", zap.Error(err))
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
	resp, err := s.rpc.ArticleClient.DeleteArticle(ctx, &articlepb.DeleteArticleRequest{Id: req.ID, AuthorID: req.AuthorID})
	if err != nil {
		s.log.Error("ArticleService/Delete error", zap.Error(err))
		return nil, err
	}
	return dto.ToArticleBoolResponse(resp.GetSuccess()), nil
}

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
	resp, err := s.rpc.ArticleClient.ListMyArticles(ctx, &articlepb.ListMyArticlesRequest{AuthorID: req.AuthorID, Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		s.log.Error("ArticleService/ListByUserID error", zap.Error(err))
		return nil, err
	}
	return dto.ToListArticlesResponse(resp.GetArticles()), nil
}

func (s *ArticleService) ListByCategory(ctx context.Context, req dto.ListByCategoryRequest) (*dto.ListArticlesResponse, error) {
	resp, err := s.rpc.ArticleClient.ListByCategory(ctx, &articlepb.ListByCategoryRequest{CategoryID: req.CategoryID, Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		s.log.Error("ArticleService/ListByCategory error", zap.Error(err))
		return nil, err
	}
	return dto.ToListArticlesResponse(resp.GetArticles()), nil
}
