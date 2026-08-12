package application

import (
	"context"
	"errors"
	"gateway/internal/client/http"
	"gateway/internal/model/dto"
	"strings"
)

type VectorService struct{ client *http.PyClient }

func NewVectorService(client *http.PyClient) *VectorService {
	return &VectorService{client: client}
}

func (s *VectorService) ListCollections(ctx context.Context) ([]dto.VectorCollection, error) {
	return s.client.ListCollections(ctx)
}

func (s *VectorService) CreateCollection(ctx context.Context, name string) (*dto.VectorCollection, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("collection name is required")
	}
	return s.client.CreateCollection(ctx, name)
}

func (s *VectorService) DeleteCollection(ctx context.Context, name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("collection name is required")
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
