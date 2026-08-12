package service

import (
	"context"
	"fmt"
	"gin-gorm-coupon-service/internal/cache"
	"gin-gorm-coupon-service/internal/pkg/crypt"
	"gin-gorm-coupon-service/internal/user/model"
	"math/rand"
	"time"
)

// UserRepo 定义 User 模块 Repository 层的接口契约
type UserRepo interface {
	CreateUser(ctx context.Context, user *model.User) error
	UserExists(ctx context.Context, username, phone string) (bool, error)
}

// UserService 业务逻辑层
type UserService struct {
	userRepo UserRepo
}

// NewUserService 构造函数
func NewUserService(userRepo UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, username, password, phone, code string) error {
	//校验验证码
	cacheCode, err := cache.GetCode(ctx, phone)
	if err != nil || cacheCode == "" {
		return fmt.Errorf("验证码不存在或已过期")
	}
	if cacheCode != code {
		return fmt.Errorf("验证码错误")
	}
	//校验用户名 / 手机号唯一
	exists, err := s.userRepo.UserExists(ctx, username, phone)
	if err != nil {
		return fmt.Errorf("查询用户失败")
	}
	if exists {
		return fmt.Errorf("用户名或手机号已存在")
	}

	//密码加密
	hashPwd, err := crypt.Hash(password)
	if err != nil {
		return fmt.Errorf("密码加密失败")
	}

	//写库
	user := &model.User{
		Username: username,
		Password: hashPwd,
		Phone:    phone,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("注册失败")
	}

	//删除验证码
	_ = cache.DelCode(ctx, phone)

	return nil
}

// SendCode 发送短信验证码
func (s *UserService) SendCode(ctx context.Context, phone string) error {
	existCode, err := cache.GetCode(ctx, phone)
	if err == nil && existCode != "" {
		// 已存在有效验证码
		return fmt.Errorf("验证码已发送，请稍后再试")
	}
	//生成 6 位验证码
	code := generateCode()
	//存入Redis
	err = cache.SetCode(ctx, phone, code, 5*time.Minute)
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
