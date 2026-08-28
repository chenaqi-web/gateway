package application

import (
	"context"
	"gateway/internal/infras/clog"

	"gateway/internal/client/rpc"
	"gateway/internal/client/rpc/core-rpc/categorypb"
	"gateway/internal/model/dto"

	"go.uber.org/zap"
)

type CategoryService struct {
	rpc *rpc.Client
	log *clog.Log
}

func NewCategoryService(rpcClient *rpc.Client, log *clog.Log) *CategoryService {
	return &CategoryService{
		rpc: rpcClient,
		log: log,
	}
}

func (s *CategoryService) CreateType(ctx context.Context, req dto.CreateTypeRequest) (*dto.CategoryBoolResponse, error) {
	resp, err := s.rpc.CategoryClient.CreateType(ctx, &categorypb.CreateTypeRequest{Name: req.Name})
	if err != nil {
		s.log.Error("CategoryService/CreateType error", zap.Error(err))
		return nil, err
	}
	return dto.ToCategoryBoolResponse(resp.GetSuccess()), nil
}

func (s *CategoryService) DeleteType(ctx context.Context, req dto.DeleteTypeRequest) (*dto.CategoryBoolResponse, error) {
	resp, err := s.rpc.CategoryClient.DeleteType(ctx, &categorypb.DeleteTypeRequest{Id: req.ID})
	if err != nil {
		s.log.Error("CategoryService/DeleteType error", zap.Error(err))
		return nil, err
	}
	return dto.ToCategoryBoolResponse(resp.GetSuccess()), nil
}

func (s *CategoryService) ListTypes(ctx context.Context) (*dto.ListTypesResponse, error) {
	resp, err := s.rpc.CategoryClient.ListTypes(ctx, &categorypb.ListTypesRequest{})
	if err != nil {
		s.log.Error("CategoryService/ListTypes error", zap.Error(err))
		return nil, err
	}
	return dto.ToListTypesResponse(resp), nil
}

func (s *CategoryService) CreateCategory(ctx context.Context, req dto.CreateCategoryRequest) (*dto.CategoryBoolResponse, error) {
	resp, err := s.rpc.CategoryClient.CreateCategory(ctx, &categorypb.CreateCategoryRequest{ParentID: req.ParentID, Name: req.Name})
	if err != nil {
		s.log.Error("CategoryService/CreateCategory error", zap.Error(err))
		return nil, err
	}
	return dto.ToCategoryBoolResponse(resp.GetSuccess()), nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, req dto.DeleteCategoryRequest) (*dto.CategoryBoolResponse, error) {
	resp, err := s.rpc.CategoryClient.DeleteCategory(ctx, &categorypb.DeleteCategoryRequest{Id: req.ID})
	if err != nil {
		s.log.Error("CategoryService/DeleteCategory error", zap.Error(err))
		return nil, err
	}
	return dto.ToCategoryBoolResponse(resp.GetSuccess()), nil
}

func (s *CategoryService) ListCategories(ctx context.Context, req dto.ListCategoriesRequest) (*dto.ListCategoriesResponse, error) {
	resp, err := s.rpc.CategoryClient.ListCategories(ctx, &categorypb.ListCategoriesRequest{ParentID: req.ParentID})
	if err != nil {
		s.log.Error("CategoryService/ListCategories error", zap.Error(err))
		return nil, err
	}
	return dto.ToListCategoriesResponse(resp), nil
}
