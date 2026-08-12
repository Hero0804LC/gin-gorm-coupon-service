package cache

import (
	"context"
	"fmt"
	"gin-gorm-coupon-service/internal/pkg/redis"
	"time"
)

const codeKeyPrefix = "code:"

func SetCode(ctx context.Context, phone, code string, ttl time.Duration) error {
	key := fmt.Sprintf("%s%s", codeKeyPrefix, phone)
	return redis.Client.Set(ctx, key, code, ttl).Err()
}

func GetCode(ctx context.Context, phone string) (string, error) {
	key := fmt.Sprintf("%s%s", codeKeyPrefix, phone)
	return redis.Client.Get(ctx, key).Result()
}

func DelCode(ctx context.Context, phone string) error {
	key := fmt.Sprintf("%s%s", codeKeyPrefix, phone)
	return redis.Client.Del(ctx, key).Err()
}
