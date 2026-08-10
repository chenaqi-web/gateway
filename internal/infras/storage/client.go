package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"gateway/internal/config"
)

type Client struct {
	impl     Provider
	provider string
}

func NewClient(cfg *config.Config) (*Client, error) {
	name := strings.ToLower(strings.TrimSpace(cfg.Storage.Provider))
	if name == "" {
		name = ProviderLocal
	}
	if name != ProviderLocal {
		return nil, fmt.Errorf("unsupported storage provider: %s", name)
	}
	impl, err := newLocalProvider(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{impl: impl, provider: name}, nil
}

func (c *Client) Provider() string { return c.provider }
func (c *Client) Upload(ctx context.Context, file *multipart.FileHeader, directory string) (*UploadResult, error) {
	return c.impl.Upload(ctx, file, directory)
}
func (c *Client) Delete(ctx context.Context, key string) error { return c.impl.Delete(ctx, key) }
func (c *Client) GetURL(key string) string                     { return c.impl.GetURL(key) }
