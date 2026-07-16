package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"backend/gateway/internal/config"

	"github.com/redis/go-redis/v9"
)

type CacheClient struct {
	Cache *redis.Client
}

// NewClient 参数设置参考
// https://aws.amazon.com/cn/blogs/china/all-roads-lead-to-rome-use-go-redis-to-connect-amazon-elasticache-for-redis-cluster/
func NewClient(cfg *config.Config) *CacheClient {
	options := redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,

		MinIdleConns: 10,

		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,

		MaxRetries:      3,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
	}

	rdb := redis.NewClient(&options)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("failed:", err)
	}

	log.Println("redis connected successfully")

	return &CacheClient{
		Cache: rdb,
	}
}

func (c *CacheClient) Close() error {
	return c.Cache.Close()
}
