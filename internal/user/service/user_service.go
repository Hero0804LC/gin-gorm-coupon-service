package service

import (
	"context"
	"fmt"

	"math/rand"
	"time"

	"gin-gorm-coupon-service/internal/cache"
)

type UserService struct {
	userRepo UserRepo
}

func NewUserService(userRepo UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

// SendCode 发送验证码
func (s *UserService) SendCode(ctx context.Context, phone string) error {
	// 1. 生成 6 位验证码
	code := generateCode()
	// 2. 存入 Redis（5分钟）
	if err := cache.SetCode(ctx, phone, code, 5*time.Minute); err != nil {
		return fmt.Errorf("保存验证码失败")
	}
	// 3. 模拟发送短信（生产环境调用短信平台）
	fmt.Printf("模拟发送短信：phone=%s, code=%s\n", phone, code)

	return nil
}

// generateCode 生成 6 位随机数字
func generateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
