package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/lagom/lagom-server/config"
	goredis "github.com/redis/go-redis/v9"
)

// Client 是 Redis 客户端封装
type Client struct {
	*goredis.Client
}

// New 创建 Redis 客户端（带连接池）
func New(cfg *config.RedisConfig) (*Client, error) {
	opt := &goredis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolTimeout:  5 * time.Second,
	}

	rdb := goredis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &Client{rdb}, nil
}

// Close 关闭连接
func (c *Client) Close() error {
	return c.Client.Close()
}
