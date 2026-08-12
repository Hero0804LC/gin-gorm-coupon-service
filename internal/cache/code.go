package cache

import (
	"context"
	"fmt"
	"time"

	"gin-gorm-coupon-service/internal/pkg/redis"
)

const codeKeyPrefix = "sms:code:"

// SetCode 保存验证码，TTL 5分钟
func SetCode(ctx context.Context, phone, code string, ttl time.Duration) error {
	key := codeKeyPrefix + phone
	return redis.Client.Set(ctx, key, code, ttl).Err()
}

// GetCode 获取验证码
func GetCode(ctx context.Context, phone string) (string, error) {
	key := codeKeyPrefix + phone
	return redis.Client.Get(ctx, key).Result()
}

// DelCode 删除验证码
func DelCode(ctx context.Context, phone string) error {
	key := fmt.Sprintf("%s%s", codeKeyPrefix, phone)
	return redis.Client.Del(ctx, key).Err()
}
