package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"gateway/internal/config"
)

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
func (c *Client) UploadAvatar(ctx context.Context, userID uint64, file *multipart.FileHeader) (*UploadResult, error) {
	return c.impl.UploadAvatar(ctx, userID, file)
}

func (c *Client) UploadCover(ctx context.Context, userID uint64, sessionID string, file *multipart.FileHeader) (*UploadResult, error) {
	return c.impl.UploadCover(ctx, userID, sessionID, file)
}

func (c *Client) UploadContent(ctx context.Context, userID uint64, sessionID string, file *multipart.FileHeader) (*UploadResult, error) {
	return c.impl.UploadContent(ctx, userID, sessionID, file)
}

func (c *Client) CommitSession(ctx context.Context, userID uint64, sessionID string) error {
	return c.impl.CommitSession(ctx, userID, sessionID)
}

func (c *Client) CleanupSessions(ctx context.Context, maxAge time.Duration) error {
	return c.impl.CleanupSessions(ctx, maxAge)
}

func (c *Client) Delete(ctx context.Context, key string) error {
	return c.impl.Delete(ctx, key)
}

func (c *Client) GetURL(key string) string {
	return c.impl.GetURL(key)
}
