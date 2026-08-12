package service

import (
	"context"
	"fmt"
	"gin-gorm-coupon-service/internal/cache"
	"math/rand"
	"time"
)

// UserRepo 定义 User 模块 Repository 层的接口契约
type UserRepo interface {
}

// UserService 业务逻辑层
type UserService struct {
	userRepo UserRepo
}

// NewUserService 构造函数
func NewUserService(userRepo UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

// SendCode 发送短信验证码
func (s *UserService) SendCode(ctx context.Context, phone string) error {
	//生成 6 位验证码
	code := generateCode()
	//存入Redis
	err := cache.SetCode(ctx, phone, code, 5*time.Minute)
	if err != nil {
		return fmt.Errorf("保存验证码失败")
	}
	//模拟发送短信
	fmt.Printf("模拟发送短信：phone=%s, code=%s\n", phone, code)
	return nil
}

// generateCode 生成 6 位随机数字
func generateCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", r.Intn(1000000))
}
