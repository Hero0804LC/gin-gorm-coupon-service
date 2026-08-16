package repository

import (
	"context"

	"gin-gorm-coupon-service/internal/order/model"
	"gorm.io/gorm"
)

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) OrderRepo {
	return &orderRepo{db: db}
}

type OrderRepo interface {
	Create(ctx context.Context, order *model.Order) error
	GetByOrderNo(ctx context.Context, orderNo string) (*model.Order, error)
	GetByUserID(ctx context.Context, userID uint64) ([]*model.Order, error)
}

func (r *orderRepo) Create(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *orderRepo) GetByOrderNo(ctx context.Context, orderNo string) (*model.Order, error) {
	var o model.Order
	if err := r.db.WithContext(ctx).Where("order_no = ?", orderNo).First(&o).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *orderRepo) GetByUserID(ctx context.Context, userID uint64) ([]*model.Order, error) {
	var list []*model.Order
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&list).Error
	return list, err
}
