package repository

import (
	"context"
	"fmt"

	orderModel "gin-gorm-coupon-service/internal/order/model"
	productModel "gin-gorm-coupon-service/internal/product/model"

	"gorm.io/gorm"
)

type OrderRepo interface {
	CreateOrder(ctx context.Context, order *orderModel.Order, items []*orderModel.OrderItem) error
	GetByOrderNo(ctx context.Context, orderNo string) (*orderModel.Order, error)
	ListByUserID(ctx context.Context, userID uint64, page, pageSize int) ([]*orderModel.Order, int64, error)
	GetByID(ctx context.Context, id, userID uint64) (*orderModel.Order, error)
	UpdateStatus(ctx context.Context, orderNo string, status int8) error
}

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) OrderRepo {
	return &orderRepo{db: db}
}

// CreateOrder 事务：创建订单主表 + 明细 + 扣库存（乐观锁）
func (r *orderRepo) CreateOrder(ctx context.Context, order *orderModel.Order, items []*orderModel.OrderItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 创建订单主表
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 2. 创建订单明细
		for _, item := range items {
			item.OrderID = order.ID
			item.OrderNo = order.OrderNo
			if err := tx.Create(item).Error; err != nil {
				return err
			}
		}

		// 3. 扣库存（乐观锁：WHERE stock >= quantity）
		for _, item := range items {
			result := tx.Model(&productModel.Product{}).
				Where("id = ? AND stock >= ?", item.ProductID, item.Quantity).
				UpdateColumn("stock", gorm.Expr("stock - ?", item.Quantity))
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return fmt.Errorf("商品 [%s] 库存不足", item.ProductName)
			}
		}

		return nil
	})
}

func (r *orderRepo) GetByOrderNo(ctx context.Context, orderNo string) (*orderModel.Order, error) {
	var order orderModel.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		Where("order_no = ?", orderNo).
		First(&order).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &order, err
}

func (r *orderRepo) ListByUserID(ctx context.Context, userID uint64, page, pageSize int) ([]*orderModel.Order, int64, error) {
	var list []*orderModel.Order
	var total int64

	query := r.db.WithContext(ctx).Model(&orderModel.Order{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Preload("Items").
		Find(&list).Error

	return list, total, err
}

func (r *orderRepo) GetByID(ctx context.Context, id, userID uint64) (*orderModel.Order, error) {
	var order orderModel.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		Where("id = ? AND user_id = ?", id, userID).
		First(&order).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &order, err
}

func (r *orderRepo) UpdateStatus(ctx context.Context, orderNo string, status int8) error {
	return r.db.WithContext(ctx).
		Model(&orderModel.Order{}).
		Where("order_no = ?", orderNo).
		Update("status", status).Error
}
