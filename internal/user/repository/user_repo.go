package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// CreateUser 创建用户
func (r *UserRepo) CreateUser(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByPhone 根据手机号查询用户
func (r *UserRepo) GetByPhone(ctx context.Context, phone string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).
		Where("phone = ?", phone).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// GetByUsername 根据用户名查询用户
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// UserExists 判断用户名或手机号是否已存在
func (r *UserRepo) UserExists(ctx context.Context, username, phone string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&User{}).
		Where("username = ? OR phone = ?", username, phone).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}
