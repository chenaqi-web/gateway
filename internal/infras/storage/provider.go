package storage

import (
	"context"
	"mime/multipart"
)

const ProviderLocal = "local"

type UploadResult struct {
	URL string
	Key string
}

type Provider interface {
	Upload(ctx context.Context, file *multipart.FileHeader, directory string) (*UploadResult, error)
	Delete(ctx context.Context, key string) error
	GetURL(key string) string
}
