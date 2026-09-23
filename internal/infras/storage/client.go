package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"gateway/internal/config"
)

const (
	ProviderLocal  = "local"
	ProviderALiYun = "aliyun"
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
	switch name {
	case ProviderLocal:
		return &Client{
			impl:     newLocalProvider(cfg),
			provider: name,
		}, nil
	case ProviderALiYun:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown storage provider: %s", cfg.Storage.Provider)
	}
}

func (c *Client) Provider() string {
	return c.provider
}

func (c *Client) Upload(ctx context.Context, file *multipart.FileHeader, directory string) (string, error) {
	return c.impl.Upload(ctx, file, directory)
}

func (c *Client) Delete(ctx context.Context, key string) error {
	return c.impl.Delete(ctx, key)
}

func (c *Client) GetURL(key string) string {
	return c.impl.GetURL(key)
}
