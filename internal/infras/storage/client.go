package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"gateway/internal/config"
)

// Client 对外暴露的存储客户端，内部按 provider 选择具体实现。
type Client struct {
	impl     Provider
	provider string
}

func NewClient(cfg *config.Config) (*Client, error) {
	providerName := strings.ToLower(strings.TrimSpace(cfg.Storage.Provider))
	if providerName == "" {
		providerName = ProviderLocal
	}

	var impl Provider
	var err error
	switch providerName {
	case ProviderLocal:
		impl, err = newLocalProvider(cfg)
	case ProviderCOS:
		return nil, fmt.Errorf("storage provider cos is not implemented yet")
	case ProviderOSS:
		return nil, fmt.Errorf("storage provider oss is not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported storage provider: %s", providerName)
	}
	if err != nil {
		return nil, err
	}

	return &Client{
		impl:     impl,
		provider: providerName,
	}, nil
}

func (c *Client) Provider() string {
	return c.provider
}

func (c *Client) Upload(ctx context.Context, file *multipart.FileHeader) (*UploadResult, error) {
	return c.impl.Upload(ctx, file)
}

func (c *Client) Delete(ctx context.Context, key string) error {
	return c.impl.Delete(ctx, key)
}

func (c *Client) GetURL(key string) string {
	return c.impl.GetURL(key)
}
