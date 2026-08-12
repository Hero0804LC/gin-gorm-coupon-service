package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var Client *redis.Client

// Init 初始化 Redis 连接
func Init(addr, password string, db int) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		PoolSize: 10,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := Client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("redis connect failed: %w", err)
	}

	return nil
}

// Close 关闭连接
func Close() {
	if Client != nil {
		Client.Close()
	}
}
