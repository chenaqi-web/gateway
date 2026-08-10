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

func (s *localProvider) Upload(_ context.Context, file *multipart.FileHeader, directory string) (*UploadResult, error) {
	now := time.Now()
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == ".jpeg" {
		ext = ".jpg"
	}
	name := fmt.Sprintf("%d%s", now.UnixNano(), ext)
	key := filepath.Join(directory, now.Format("2006"), now.Format("01"), now.Format("02"), name)
	path := filepath.Join(s.basePath, key)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	source, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer source.Close()
	target, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer target.Close()
	if _, err := io.Copy(target, source); err != nil {
		return nil, err
	}
	urlKey := s.urlPrefix + "/" + filepath.ToSlash(key)
	return &UploadResult{URL: s.baseURL + urlKey, Key: urlKey}, nil
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

func (s *localProvider) GetURL(key string) string {
	if strings.HasPrefix(key, "http://") || strings.HasPrefix(key, "https://") {
		return key
	}
	return s.baseURL + "/" + strings.TrimLeft(key, "/")
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
