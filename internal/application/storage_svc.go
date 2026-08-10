package application

import (
	"context"
	"errors"
	"gateway/internal/utils"
	"mime/multipart"
	"strconv"
	"strings"
	"sync"
	"time"

	"gateway/internal/config"
	"gateway/internal/infras/storage"
)

const defaultMaxUploadSize = 10 << 20
const abandonedSessionAge = 24 * time.Hour

var (
	ErrStorageMissingFile  = errors.New("file is required")
	ErrStorageFileTooLarge = errors.New("file exceeds max upload size")
	ErrStorageInvalidImage = errors.New("only jpg, png, gif, webp and avif images are allowed")
	ErrStorageInvalidKey   = errors.New("invalid storage key")
)

type StorageService struct {
	cfg         *config.Config
	client      *storage.Client
	cleanupMu   sync.Mutex
	lastCleanup time.Time
}

type UploadResponse struct {
	URL      string `json:"url"`
	Key      string `json:"key"`
	Provider string `json:"provider"`
}

func NewStorageService(cfg *config.Config, client *storage.Client) *StorageService {
	return &StorageService{
		cfg:    cfg,
		client: client,
	}
}

func (s *StorageService) UploadAvatar(ctx context.Context, userID uint64, file *multipart.FileHeader) (*UploadResponse, error) {
	if err := utils.ValidateImage(s.cfg, file); err != nil {
		return nil, err
	}
	result, err := s.client.UploadAvatar(ctx, userID, file)
	if err != nil {
		return nil, err
	}
	return &UploadResponse{URL: result.URL, Key: result.Key, Provider: s.client.Provider()}, nil
}

func (s *StorageService) UploadCover(ctx context.Context, userID uint64, sessionID string, file *multipart.FileHeader) (*UploadResponse, error) {
	if err := utils.ValidateSession(sessionID); err != nil {
		return nil, err
	}

	if err := utils.ValidateImage(s.cfg, file); err != nil {
		return nil, err
	}

	s.cleanupExpired(ctx)
	result, err := s.client.UploadCover(ctx, userID, sessionID, file)
	if err != nil {
		return nil, err
	}
	return &UploadResponse{URL: result.URL, Key: result.Key, Provider: s.client.Provider()}, nil
}

func (s *StorageService) UploadContent(ctx context.Context, userID uint64, sessionID string, file *multipart.FileHeader) (*UploadResponse, error) {
	if err := utils.ValidateSession(sessionID); err != nil {
		return nil, err
	}

	if err := utils.ValidateImage(s.cfg, file); err != nil {
		return nil, err
	}

	s.cleanupExpired(ctx)
	result, err := s.client.UploadContent(ctx, userID, sessionID, file)
	if err != nil {
		return nil, err
	}
	return &UploadResponse{URL: result.URL, Key: result.Key, Provider: s.client.Provider()}, nil
}

func (s *StorageService) CommitSession(ctx context.Context, userID uint64, sessionID string) error {
	if err := utils.ValidateSession(sessionID); err != nil {
		return err
	}
	return s.client.CommitSession(ctx, userID, sessionID)
}

func (s *StorageService) DeleteOwned(ctx context.Context, userID uint64, key string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return ErrStorageInvalidKey
	}
	ownedBlogPrefix := "/static/upload/blog/" + strconv.FormatUint(userID, 10) + "/"
	ownedAvatarPrefix := "/static/upload/avatar/" + strconv.FormatUint(userID, 10) + "."
	if !strings.Contains(key, ownedBlogPrefix) && !strings.Contains(key, ownedAvatarPrefix) {
		return ErrStorageInvalidKey
	}
	return s.client.Delete(ctx, key)
}

func (s *StorageService) cleanupExpired(ctx context.Context) {
	s.cleanupMu.Lock()
	defer s.cleanupMu.Unlock()
	if time.Since(s.lastCleanup) < time.Hour {
		return
	}
	if s.client.CleanupSessions(ctx, abandonedSessionAge) == nil {
		s.lastCleanup = time.Now()
	}
}
