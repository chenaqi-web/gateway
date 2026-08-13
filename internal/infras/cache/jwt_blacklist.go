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
	return fmt.Sprintf("auth:blacklist:token:%x", sha256.Sum256([]byte(token)))
}

// =====================================================================================================================
// 用户被强制拉黑时的缓存

func (c *JwtBlacklist) AddBlacklist(ctx context.Context, userID uint64) error {
	return c.Cache.Set(ctx, userIDBlacklistKey(userID), "", 0).Err()
}

func (c *JwtBlacklist) RemoveBlacklist(ctx context.Context, userID uint64) error {
	return c.Cache.Del(ctx, userIDBlacklistKey(userID)).Err()
}

func (c *JwtBlacklist) IsBlacklisted(ctx context.Context, userID uint64) (bool, error) {
	exists, err := c.Cache.Exists(ctx, userIDBlacklistKey(userID)).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func userIDBlacklistKey(userID uint64) string {
	return fmt.Sprintf("auth:blacklist:ID:%d", userID)
}
