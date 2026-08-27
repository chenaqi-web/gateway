package cache

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"
)

type Blacklist struct {
	*CacheClient
}

func NewJwtBlacklist(client *CacheClient) *Blacklist {
	return &Blacklist{client}
}

func (c *Blacklist) AddToken(ctx context.Context, token string, expireSeconds int) error {
	return c.Cache.Set(
		ctx,
		tokenBlacklistKey(token),
		"",
		time.Duration(expireSeconds)*time.Second,
	).Err()
}

func (c *Blacklist) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	count, err := c.Cache.Exists(ctx, tokenBlacklistKey(token)).Result()
	return count > 0, err
}

func tokenBlacklistKey(token string) string {
	return fmt.Sprintf("blacklist:token:%x", sha256.Sum256([]byte(token)))
}

// =====================================================================================================================
// 用户被强制拉黑时的缓存

func (c *Blacklist) AddUser(ctx context.Context, userID uint64) error {
	return c.Cache.Set(ctx, userIDBlacklistKey(userID), "", 0).Err()
}

func (c *Blacklist) RemoveUser(ctx context.Context, userID uint64) error {
	return c.Cache.Del(ctx, userIDBlacklistKey(userID)).Err()
}

func (c *Blacklist) IsUserBlacklisted(ctx context.Context, userID uint64) (bool, error) {
	exists, err := c.Cache.Exists(ctx, userIDBlacklistKey(userID)).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func userIDBlacklistKey(userID uint64) string {
	return fmt.Sprintf("auth:blacklist:ID:%d", userID)
}
