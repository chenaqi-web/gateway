package storage

import (
	"context"
	"mime/multipart"
)

const (
	ProviderLocal = "local"
	ProviderCOS   = "cos"
	ProviderOSS   = "oss"
)

// UploadResult 上传结果。
type UploadResult struct {
	URL string // 完整访问 URL
	Key string // 对象 key / 相对路径
}

// Provider 存储后端接口，local / cos / oss 均实现该接口。
type Provider interface {
	Upload(ctx context.Context, file *multipart.FileHeader) (*UploadResult, error)
	Delete(ctx context.Context, key string) error
	GetURL(key string) string
}
