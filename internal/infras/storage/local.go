package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gateway/internal/config"
)

type localProvider struct {
	basePath  string // 静态文件的存储位置
	baseURL   string // 访问的url前缀（域名等）
	urlPrefix string //
}

const localURLPrefix = "/static/upload"

func newLocalProvider(cfg *config.Config) Provider {
	return &localProvider{
		basePath:  cfg.Storage.BasePath,
		baseURL:   cfg.Storage.BaseURL,
		urlPrefix: localURLPrefix,
	}
}

func (s *localProvider) Upload(_ context.Context, file *multipart.FileHeader, directory string) (string, error) {
	// 1.构建存储的文件名 时间戳.ext
	now := time.Now()
	ext := strings.ToLower(filepath.Ext(file.Filename))
	name := fmt.Sprintf("%d%s", now.UnixNano(), ext)
	key := filepath.Join(directory, name)
	path := filepath.Join(s.basePath, key)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	source, err := file.Open()
	if err != nil {
		return "", err
	}
	defer source.Close()
	target, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer target.Close()
	if _, err := io.Copy(target, source); err != nil {
		return "", err
	}
	urlKey := s.urlPrefix + "/" + filepath.ToSlash(key)
	return s.baseURL + urlKey, nil
}

func (s *localProvider) GetURL(key string) string {
	if strings.HasPrefix(key, "http://") || strings.HasPrefix(key, "https://") {
		return key
	}
	return s.baseURL + "/" + strings.TrimLeft(key, "/")
}

func (s *localProvider) Delete(_ context.Context, key string) error {
	path, err := s.toLocalPath(key)
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (s *localProvider) toLocalPath(key string) (string, error) {
	key = strings.TrimPrefix(key, s.baseURL)
	key = strings.TrimPrefix(key, s.urlPrefix)
	clean := filepath.Clean(strings.TrimLeft(key, "/\\"))
	if clean == "." || strings.HasPrefix(clean, "..") {
		return "", fmt.Errorf("invalid storage key")
	}
	path := filepath.Join(s.basePath, clean)
	base, _ := filepath.Abs(s.basePath)
	absolute, _ := filepath.Abs(path)
	if absolute != base && !strings.HasPrefix(absolute, base+string(os.PathSeparator)) {
		return "", fmt.Errorf("invalid storage key")
	}
	return path, nil
}
