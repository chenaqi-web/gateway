package storage

import (
	"context"
	"errors"
	"mime/multipart"
)

var ErrInvalidStorageKey = errors.New("invalid storage key")

// 业务目录建议只使用这些值。
const (
	DirectoryAvatar         = "avatar"
	DirectoryArticleCover   = "article/cover"
	DirectoryArticleContent = "article/content"
)

type Provider interface {
	Upload(ctx context.Context, file *multipart.FileHeader, directory string, id uint64) (string, error)
	Delete(ctx context.Context, key string) error
}
