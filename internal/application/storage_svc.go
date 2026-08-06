package application

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"

	"gateway/internal/config"
	"gateway/internal/infras/storage"
)

const defaultMaxUploadSize = 10 << 20 // 10MB

var (
	ErrStorageMissingFile  = errors.New("file is required")
	ErrStorageFileTooLarge = errors.New("file exceeds max upload size")
	ErrStorageInvalidKey   = errors.New("storage key is required")
)

type StorageService struct {
	cfg    *config.Config
	client *storage.Client
}

func NewStorageService(cfg *config.Config, client *storage.Client) *StorageService {
	return &StorageService{
		cfg:    cfg,
		client: client,
	}
}

type UploadResponse struct {
	URL      string
	Key      string
	Provider string
}

func (s *StorageService) Upload(ctx context.Context, file *multipart.FileHeader) (*UploadResponse, error) {
	if file == nil {
		return nil, ErrStorageMissingFile
	}
	if file.Size > s.maxUploadSize() {
		return nil, fmt.Errorf("%w: max %d bytes", ErrStorageFileTooLarge, s.maxUploadSize())
	}

	result, err := s.client.Upload(ctx, file)
	if err != nil {
		return nil, err
	}
	return &UploadResponse{URL: result.URL, Key: result.Key, Provider: s.client.Provider()}, nil
}

func (s *StorageService) Delete(ctx context.Context, key string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return ErrStorageInvalidKey
	}
	return s.client.Delete(ctx, key)
}

func (s *StorageService) maxUploadSize() int64 {
	if s.cfg.Storage.MaxSize > 0 {
		return s.cfg.Storage.MaxSize
	}
	return defaultMaxUploadSize
}
