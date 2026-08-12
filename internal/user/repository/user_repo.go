package repository

import (
	"context"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// 创建用户
func (r *UserRepo) CreateUser(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// 根据用户名或手机号查用户
func (r *UserRepo) GetByUsernameOrPhone(ctx context.Context, username, phone string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).
		Where("username = ? OR phone = ?", username, phone).
		First(&user).Error
	return &user, err
}

// 根据 ID 查用户
func (r *UserRepo) GetByID(ctx context.Context, id uint64) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).First(&user, id).Error
	return &user, err
}
