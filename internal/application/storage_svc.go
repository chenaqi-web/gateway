package application

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"

	"backend/gateway/internal/config"
	"backend/gateway/internal/infras/storage"
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

func (s *StorageService) Upload(ctx context.Context, file *multipart.FileHeader) (*storage.UploadResult, error) {
	if file == nil {
		return nil, ErrStorageMissingFile
	}
	if file.Size > s.maxUploadSize() {
		return nil, fmt.Errorf("%w: max %d bytes", ErrStorageFileTooLarge, s.maxUploadSize())
	}

	return s.client.Upload(ctx, file)
}

func (s *StorageService) Delete(ctx context.Context, key string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return ErrStorageInvalidKey
	}
	return s.client.Delete(ctx, key)
}

func (s *StorageService) GetURL(key string) string {
	return s.client.GetURL(key)
}

func (s *StorageService) Provider() string {
	return s.client.Provider()
}

func (s *StorageService) maxUploadSize() int64 {
	if s.cfg.Storage.MaxSize > 0 {
		return s.cfg.Storage.MaxSize
	}
	return defaultMaxUploadSize
}
