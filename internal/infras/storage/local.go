package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gateway/internal/config"
)

type localProvider struct {
	BaseUrl   string
	urlPrefix string // 对外URL前缀
	rootDir   string // 静态文件的存储位置
}

const localURLPrefix = "/static/upload/"

func newLocalProvider(cfg *config.Config) Provider {
	return &localProvider{
		cfg.Storage.BaseURL,
		localURLPrefix,
		cfg.Storage.RootDir,
	}
}

func (s *localProvider) Upload(_ context.Context, file *multipart.FileHeader, directory string, id uint64) (string, error) {
	if filepath.IsAbs(directory) {
		return "", ErrInvalidStorageKey
	}
	directory = filepath.ToSlash(filepath.Clean(directory))
	if directory == "." || strings.HasPrefix(directory, "../") || directory == ".." {
		return "", ErrInvalidStorageKey
	}

	now := time.Now().UTC()
	ext := strings.ToLower(filepath.Ext(filepath.Base(file.Filename)))
	key := filepath.ToSlash(filepath.Join(directory, now.Format("2006/01"), fmt.Sprintf("%d-%d%s", now.UnixMilli(), id, ext)))

	key, err := s.write(file, key)
	if err != nil {
		return "", err
	}

	return s.BaseUrl + s.urlPrefix + key, nil
}

func (s *localProvider) write(file *multipart.FileHeader, key string) (string, error) {
	path := filepath.Join(s.rootDir, filepath.FromSlash(key))

	// 递归创建目录
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}

	// 打开上传文件流
	source, err := file.Open()
	if err != nil {
		return "", err
	}
	defer source.Close()

	// 创建临时文件，前缀是.upload-
	tmp, err := os.CreateTemp(filepath.Dir(path), ".upload-*")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	// 将文件拷贝到临时文件
	if _, err = io.Copy(tmp, source); err != nil {
		_ = tmp.Close()
		return "", err
	}
	// 关闭临时文件，缓冲要flush到磁盘，否则无法呗Rename
	if err = tmp.Close(); err != nil {
		return "", err
	}
	if err = os.Chmod(tmpPath, 0o644); err != nil {
		return "", err
	}

	// 原子重命名：临时文件 → 目标文件
	// 这样做是为了避免看到半个文件的情况
	if err = os.Rename(tmpPath, path); err != nil {
		return "", err
	}
	return filepath.ToSlash(key), nil
}

func (s *localProvider) Delete(_ context.Context, key string) error {
	// 1.如果可以是一个绝对的URL，则提取它的Path部分
	if parsed, err := url.Parse(key); err == nil && parsed.IsAbs() {
		key = parsed.Path
	}

	// 2.去掉前缀
	key = strings.TrimPrefix(key, s.urlPrefix)
	key = strings.TrimPrefix(key, "/")
	key = filepath.ToSlash(filepath.Clean(key))
	if key == "." || strings.HasPrefix(key, "../") || key == ".." {
		return ErrInvalidStorageKey
	}

	// 3.删除照片
	if err := os.Remove(filepath.Join(s.rootDir, filepath.FromSlash(key))); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
