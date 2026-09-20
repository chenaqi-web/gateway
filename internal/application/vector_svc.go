package application

import (
	"context"
	"gateway/internal/client/http"
	"gateway/internal/infras/clog"
	"gateway/internal/model/dto"
	"strings"

	"go.uber.org/zap"
)

type VectorService struct {
	client *http.PyClient
	log    *clog.Log
}

func NewVectorService(client *http.PyClient, log *clog.Log) *VectorService {
	return &VectorService{
		client: client,
		log:    log,
	}
}

func (s *VectorService) CreateCollection(ctx context.Context, request *dto.CreateVectorCollectionRequest) error {
	if strings.TrimSpace(request.Name) == "" {
		return EmptyValueError
	}

	err := s.client.CreateCollection(ctx, request.Name)
	if err != nil {
		s.log.Error("VectorService/CreateCollection error", zap.Error(err))
		return err
	}

	return nil
}

func (s *VectorService) DeleteCollection(ctx context.Context, request *dto.DelVectorCollectionRequest) error {
	if strings.TrimSpace(request.Name) == "" {
		return EmptyValueError
	}

	if err := s.client.DeleteCollection(ctx, request.Name); err != nil {
		s.log.Error("VectorService/DeleteCollection error", zap.Error(err))
		return err
	}

	return nil
}

func (s *VectorService) ListCollections(ctx context.Context) ([]*dto.Collection, error) {
	res, err := s.client.ListCollections(ctx)
	if err != nil {
		s.log.Error("VectorService/ListCollections error", zap.Error(err))
		return nil, err
	}
	return res, nil
}

func (s *VectorService) ListDocuments(ctx context.Context, request *dto.ListDocumentsRequest) (*dto.ListDocumentsResponse, error) {
	documents, err := s.client.ListDocuments(ctx, request.Name, request.Page, request.PageSize)
	if err != nil {
		s.log.Error("VectorService/ListDocuments error", zap.Error(err))
		return nil, err
	}
	return documents, nil
}

func (s *VectorService) SearchDocuments(ctx context.Context, request *dto.DocsSearchRequest) (*dto.DocsSearchResponse, error) {
	docs, err := s.client.SearchVectors(ctx, "", request)
	if err != nil {
		s.log.Error("VectorService/SearchDocuments error", zap.Error(err))
		return nil, err
	}
	return docs, nil
}
