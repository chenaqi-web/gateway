package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
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

// todo 直接保存就可以，头像不存在临时一说

func (s *localProvider) UploadAvatar(_ context.Context, userID uint64, file *multipart.FileHeader) (*UploadResult, error) {
	//folder := "avatar"
	return nil, nil
}

func (s *localProvider) UploadCover(_ context.Context, userID uint64, sessionID string, file *multipart.FileHeader) (*UploadResult, error) {
	folder := filepath.Join("blog", strconv.FormatUint(userID, 10), sessionID)
	if err := s.touchSession(folder); err != nil {
		return nil, err
	}
	return s.saveFixed(file, folder, "cover")
}

func (s *localProvider) UploadContent(_ context.Context, userID uint64, sessionID string, file *multipart.FileHeader) (*UploadResult, error) {
	folder := filepath.Join("blog", strconv.FormatUint(userID, 10), sessionID)
	if err := s.touchSession(folder); err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%d%s", time.Now().UnixNano(), normalizedExt(file.Filename))
	return s.save(file, filepath.Join(folder, "content", name))
}

func (s *localProvider) CommitSession(_ context.Context, userID uint64, sessionID string) error {
	folder := filepath.Join(s.basePath, filepath.Join("blog", strconv.FormatUint(userID, 10), sessionID))
	if _, err := os.Stat(folder); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(folder, ".committed"), []byte(time.Now().Format(time.RFC3339)), 0o644)
}

func (s *localProvider) CleanupSessions(_ context.Context, maxAge time.Duration) error {
	root := filepath.Join(s.basePath, "blog")
	users, err := os.ReadDir(root)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	cutoff := time.Now().Add(-maxAge)
	for _, user := range users {
		if !user.IsDir() {
			continue
		}
		sessions, readErr := os.ReadDir(filepath.Join(root, user.Name()))
		if readErr != nil {
			continue
		}
		for _, session := range sessions {
			if !session.IsDir() {
				continue
			}
			dir := filepath.Join(root, user.Name(), session.Name())
			if _, markerErr := os.Stat(filepath.Join(dir, ".committed")); markerErr == nil {
				continue
			}
			info, infoErr := os.Stat(filepath.Join(dir, ".session"))
			if infoErr == nil && info.ModTime().Before(cutoff) {
				_ = os.RemoveAll(dir)
			}
		}
	}
	return nil
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

func (s *localProvider) save(file *multipart.FileHeader, objectKey string) (*UploadResult, error) {
	path := filepath.Join(s.basePath, objectKey)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()
	dst, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return nil, err
	}
	key := s.urlPrefix + "/" + filepath.ToSlash(objectKey)
	return &UploadResult{URL: s.baseURL + key, Key: key}, nil
}

func (s *localProvider) saveFixed(file *multipart.FileHeader, folder, name string) (*UploadResult, error) {
	dir := filepath.Join(s.basePath, folder)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	temporary, err := os.CreateTemp(dir, ".upload-*")
	if err != nil {
		return nil, err
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)
	source, err := file.Open()
	if err != nil {
		temporary.Close()
		return nil, err
	}
	_, copyErr := io.Copy(temporary, source)
	source.Close()
	closeErr := temporary.Close()
	if copyErr != nil {
		return nil, copyErr
	}
	if closeErr != nil {
		return nil, closeErr
	}
	if err := s.removeVariants(folder, name); err != nil {
		return nil, err
	}
	objectKey := filepath.Join(folder, name+normalizedExt(file.Filename))
	if err := os.Rename(temporaryPath, filepath.Join(s.basePath, objectKey)); err != nil {
		return nil, err
	}
	key := s.urlPrefix + "/" + filepath.ToSlash(objectKey)
	return &UploadResult{URL: s.baseURL + key, Key: key}, nil
}

func (s *localProvider) touchSession(folder string) error {
	dir := filepath.Join(s.basePath, folder)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	now := time.Now()
	marker := filepath.Join(dir, ".session")
	if err := os.WriteFile(marker, []byte(now.Format(time.RFC3339)), 0o644); err != nil {
		return err
	}
	return os.Chtimes(marker, now, now)
}

func (s *localProvider) removeVariants(folder, name string) error {
	matches, err := filepath.Glob(filepath.Join(s.basePath, folder, name+".*"))
	if err != nil {
		return err
	}
	for _, match := range matches {
		if err := os.Remove(match); err != nil && !os.IsNotExist(err) {
			return err
		}
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

func normalizedExt(name string) string {
	ext := strings.ToLower(filepath.Ext(name))
	if ext == ".jpeg" {
		return ".jpg"
	}
	return ext
}
