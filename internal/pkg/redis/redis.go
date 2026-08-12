package redis

import (
	"context"
	"fmt"
	"time"

	"gin-gorm-coupon-service/config"

	"github.com/redis/go-redis/v9"
)

// Client 全局 Redis 客户端
var Client *redis.Client

// Init 初始化 Redis
func Init() error {
	cfg := config.Get().Redis

	Client = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}

	return nil
}

// Close 关闭 Redis 连接
func Close() {
	if Client != nil {
		if err := Client.Close(); err != nil {
			fmt.Println("redis close error:", err)
		}
	}
}
