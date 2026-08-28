package application

import (
	"context"
	"errors"
	"gateway/internal/client/http"
	"gateway/internal/infras/clog"
	"gateway/internal/model/dto"
	"go.uber.org/zap"
	"strings"
)

type VectorService struct {
	client *http.PyClient
	log    *clog.Log
}

func NewVectorService(client *http.PyClient, log *clog.Log) *VectorService {
	return &VectorService{client: client, log: log}
}

func (s *VectorService) ListCollections(ctx context.Context) ([]dto.VectorCollection, error) {
	return s.client.ListCollections(ctx)
}

func (s *VectorService) CreateCollection(ctx context.Context, name string) (*dto.VectorCollection, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		err := errors.New("collection name is required")
		s.log.Error("VectorService/CreateCollection error", zap.Error(err))
		return nil, err
	}
	return s.client.CreateCollection(ctx, name)
}

func (s *VectorService) DeleteCollection(ctx context.Context, name string) error {
	if strings.TrimSpace(name) == "" {
		err := errors.New("collection name is required")
		s.log.Error("VectorService/DeleteCollection error", zap.Error(err))
		return err
	}
	return s.client.DeleteCollection(ctx, name)
}

func (s *VectorService) ListDocuments(ctx context.Context, name string, page, pageSize int) (*dto.VectorDocumentPage, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 5
	}
	if pageSize > 5 {
		pageSize = 5
	}
	return s.client.ListDocuments(ctx, name, page, pageSize)
}
