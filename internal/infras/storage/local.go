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
	basePath  string
	baseURL   string
	urlPrefix string
}

const (
	defaultLocalBasePath = "./static/upload"
	localURLPrefix       = "/static/upload"
)

func newLocalProvider(cfg *config.Config) (Provider, error) {
	basePath := strings.TrimSpace(cfg.Storage.BasePath)
	if basePath == "" {
		basePath = defaultLocalBasePath
	}

	baseURL := strings.TrimRight(strings.TrimSpace(cfg.Storage.BaseURL), "/")
	if baseURL == "" {
		baseURL = "http://" + strings.TrimSpace(cfg.Server.Addr)
	}

	if err := os.MkdirAll(basePath, 0o755); err != nil {
		return nil, fmt.Errorf("create local storage dir: %w", err)
	}

	return &localProvider{
		basePath:  basePath,
		baseURL:   baseURL,
		urlPrefix: localURLPrefix,
	}, nil
}

func (s *localProvider) Upload(_ context.Context, file *multipart.FileHeader) (*UploadResult, error) {
	objectKey := generateObjectKey(file.Filename)
	dir := filepath.Join(s.basePath, filepath.Dir(objectKey))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	savePath := filepath.Join(s.basePath, objectKey)

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.Create(savePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return nil, err
	}

	urlPath := s.urlPrefix + "/" + objectKey
	return &UploadResult{
		URL: s.baseURL + urlPath,
		Key: urlPath,
	}, nil
}

func (s *localProvider) Delete(_ context.Context, key string) error {
	localPath, err := s.toLocalPath(key)
	if err != nil {
		return err
	}
	if err := os.Remove(localPath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (s *localProvider) GetURL(key string) string {
	if strings.HasPrefix(key, "http://") || strings.HasPrefix(key, "https://") {
		return key
	}
	if !strings.HasPrefix(key, "/") {
		key = "/" + key
	}
	return s.baseURL + key
}

func (s *localProvider) toLocalPath(key string) (string, error) {
	key = strings.TrimPrefix(key, s.baseURL)
	key = strings.TrimPrefix(key, s.urlPrefix+"/")
	key = strings.TrimPrefix(key, s.urlPrefix)
	key = strings.TrimPrefix(key, "/")
	if key == "" || strings.Contains(key, "..") {
		return "", fmt.Errorf("invalid storage key: %s", key)
	}
	return filepath.Join(s.basePath, key), nil
}

func (s *localProvider) BasePath() string {
	return s.basePath
}

func generateObjectKey(filename string) string {
	now := time.Now()
	ext := filepath.Ext(filename)
	name := fmt.Sprintf("%d%s", now.UnixNano(), ext)
	return fmt.Sprintf("%s/%s/%s/%s", now.Format("2006"), now.Format("01"), now.Format("02"), name)
}
