package repository

import (
	"context"
	"time"

	"gin-gorm-coupon-service/internal/coupon/model"
	"gorm.io/gorm"
)

type couponRepo struct {
	db *gorm.DB
}

func NewCouponRepo(db *gorm.DB) CouponRepo {
	return &couponRepo{db: db}
}

type CouponRepo interface {
	Create(ctx context.Context, coupon *model.Coupon) error
	GetByID(ctx context.Context, id int64) (*model.Coupon, error)
	List(ctx context.Context) ([]*model.Coupon, error)
}

func (r *couponRepo) Create(ctx context.Context, coupon *model.Coupon) error {
	return r.db.WithContext(ctx).Create(coupon).Error
}

func (r *couponRepo) GetByID(ctx context.Context, id int64) (*model.Coupon, error) {
	var c model.Coupon
	if err := r.db.WithContext(ctx).First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *couponRepo) List(ctx context.Context) ([]*model.Coupon, error) {
	var list []*model.Coupon
	err := r.db.WithContext(ctx).
		Where("status = ? AND end_time > ?", 1, time.Now()).
		Order("created_at DESC").
		Find(&list).Error
	return list, err
}
