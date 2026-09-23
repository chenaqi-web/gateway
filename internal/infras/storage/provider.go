package storage

import (
	"context"
	"errors"
	"mime/multipart"
)

var ErrInvalidStorageKey = errors.New("invalid storage key")

// 业务目录建议只使用这些值；如果需要按用户隔离，可传入
// "avatar/user-<id>" 这样的子目录。
const (
	DirectoryAvatar         = "avatar"
	DirectoryArticleCover   = "article/cover"
	DirectoryArticleContent = "article/content"
)

type Provider interface {
	Upload(ctx context.Context, file *multipart.FileHeader, directory string) (string, error)
	UploadAvatar(ctx context.Context, file *multipart.FileHeader, userID uint64) (string, error)
	Delete(ctx context.Context, key string) error
}
