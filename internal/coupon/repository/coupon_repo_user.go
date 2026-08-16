package repository

import (
	"context"

	"gin-gorm-coupon-service/internal/coupon/model"
	"gorm.io/gorm"
)

type userCouponRepo struct {
	db *gorm.DB
}

func NewUserCouponRepo(db *gorm.DB) UserCouponRepo {
	return &userCouponRepo{db: db}
}

type UserCouponRepo interface {
	Create(ctx context.Context, uc *model.UserCoupon) error
}

func (r *userCouponRepo) Create(ctx context.Context, uc *model.UserCoupon) error {
	return r.db.WithContext(ctx).Create(uc).Error
}
