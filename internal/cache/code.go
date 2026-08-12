package cache

import (
	"context"
	"fmt"
	"time"

	"gin-gorm-coupon-service/internal/pkg/redis"
)

// 存验证码，5分钟过期
func SetCode(phone, code string) error {
	ctx := context.Background()
	key := fmt.Sprintf("code:%s", phone)
	return redis.Client.Set(ctx, key, code, time.Minute*5).Err()
}

// 取验证码
func GetCode(phone string) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("code:%s", phone)
	return redis.Client.Get(ctx, key).Result()
}
