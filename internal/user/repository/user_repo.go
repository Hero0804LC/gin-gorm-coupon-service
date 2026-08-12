package repository

import (
	"context"

	"gin-gorm-coupon-service/internal/user/model"
	"gorm.io/gorm"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *model.User) error
	UserExists(ctx context.Context, username, phone string) (bool, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}

// CreateUser 创建用户
func (r *userRepo) CreateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// UserExists 判断用户名或手机号是否存在
func (r *userRepo) UserExists(ctx context.Context, username, phone string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("username = ? OR phone = ?", username, phone).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}
