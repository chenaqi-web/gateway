package cache

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"
)

type JwtBlacklist struct {
	*CacheClient
}

func NewJwtBlacklist(client *CacheClient) *JwtBlacklist {
	return &JwtBlacklist{client}
}

func (c *JwtBlacklist) BlacklistToken(ctx context.Context, token string, expireSeconds int) error {
	return c.Cache.Set(
		ctx,
		tokenBlacklistKey(token),
		1,
		time.Duration(expireSeconds)*time.Second,
	).Err()
}

func (c *JwtBlacklist) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	count, err := c.Cache.Exists(ctx, tokenBlacklistKey(token)).Result()
	return count > 0, err
}

func tokenBlacklistKey(token string) string {
	return fmt.Sprintf("auth:blacklist:%x", sha256.Sum256([]byte(token)))
}
