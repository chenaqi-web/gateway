package storage

import (
	"context"
	"mime/multipart"
	"time"
)

const (
	ProviderLocal = "local"
	ProviderCOS   = "cos"
	ProviderOSS   = "oss"
)

type UploadResult struct {
	URL string
	Key string
}

type Provider interface {
	UploadAvatar(ctx context.Context, userID uint64, file *multipart.FileHeader) (*UploadResult, error)
	UploadCover(ctx context.Context, userID uint64, sessionID string, file *multipart.FileHeader) (*UploadResult, error)
	UploadContent(ctx context.Context, userID uint64, sessionID string, file *multipart.FileHeader) (*UploadResult, error)
	CommitSession(ctx context.Context, userID uint64, sessionID string) error
	CleanupSessions(ctx context.Context, maxAge time.Duration) error
	Delete(ctx context.Context, key string) error
	GetURL(key string) string
}
