package application

import (
	"context"
	"errors"
	"gateway/internal/utils"
	"mime/multipart"
	"strings"

	"gateway/internal/config"
	"gateway/internal/infras/storage"
)

const defaultMaxUploadSize = 10 << 20

var (
	ErrStorageMissingFile  = errors.New("file is required")
	ErrStorageFileTooLarge = errors.New("file exceeds max upload size")
	ErrStorageInvalidImage = errors.New("only jpg, png, gif, webp and avif images are allowed")
	ErrStorageInvalidKey   = errors.New("invalid storage key")
)

type StorageService struct {
	cfg    *config.Config
	client *storage.Client
}
type UploadResponse struct {
	URL      string `json:"url"`
	Key      string `json:"key"`
	Provider string `json:"provider"`
}

func NewStorageService(cfg *config.Config, client *storage.Client) *StorageService {
	return &StorageService{cfg: cfg, client: client}
}
func (s *StorageService) UploadAvatar(ctx context.Context, file *multipart.FileHeader) (*UploadResponse, error) {
	return s.uploadImage(ctx, file, "avatar")
}
func (s *StorageService) UploadCover(ctx context.Context, file *multipart.FileHeader) (*UploadResponse, error) {
	return s.uploadImage(ctx, file, "cover")
}
func (s *StorageService) UploadContent(ctx context.Context, file *multipart.FileHeader) (*UploadResponse, error) {
	return s.uploadImage(ctx, file, "content")
}

func (s *StorageService) uploadImage(ctx context.Context, file *multipart.FileHeader, directory string) (*UploadResponse, error) {
	if err := utils.ValidateImage(s.cfg, file); err != nil {
		return nil, err
	}
	result, err := s.client.Upload(ctx, file, directory)
	if err != nil {
		return nil, err
	}
	return &UploadResponse{URL: result.URL, Key: result.Key, Provider: s.client.Provider()}, nil
}

func (s *StorageService) Delete(ctx context.Context, key string) error {
	if strings.TrimSpace(key) == "" {
		return ErrStorageInvalidKey
	}
	return s.client.Delete(ctx, key)
}
