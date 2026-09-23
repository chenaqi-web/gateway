package storage

import (
	"context"
	"mime/multipart"
)

type Provider interface {
	Provider() string
	Upload(ctx context.Context, file *multipart.FileHeader, directory string) (string, error)
	Delete(ctx context.Context, key string) error
	GetURL(key string) string
}
