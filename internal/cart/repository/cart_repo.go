package repository

import (
	"context"
	"errors"

	cartModel "gin-gorm-coupon-service/internal/cart/model"
	productModel "gin-gorm-coupon-service/internal/product/model"

	"gorm.io/gorm"
)

type CartRepo interface {
	AddOrUpdate(ctx context.Context, item *cartModel.CartItem) error
	ListByUserID(ctx context.Context, userID uint64) ([]*cartModel.CartItem, error)
	UpdateQuantity(ctx context.Context, id, userID uint64, quantity int) error
	DeleteByID(ctx context.Context, id, userID uint64) error
	ClearByUserID(ctx context.Context, userID uint64) error
	GetByID(ctx context.Context, id, userID uint64) (*cartModel.CartItem, error)
	GetByUserAndProduct(ctx context.Context, userID, productID uint64) (*cartModel.CartItem, error)
}

type cartRepo struct {
	db *gorm.DB
}

func NewCartRepo(db *gorm.DB) CartRepo {
	return &cartRepo{db: db}
}

// AddOrUpdate 加入购物车：已存在则数量累加，不存在则新增
func (r *cartRepo) AddOrUpdate(ctx context.Context, item *cartModel.CartItem) error {
	var existing cartModel.CartItem
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND product_id = ?", item.UserID, item.ProductID).
		First(&existing).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 不存在 → 新增
		return r.db.WithContext(ctx).Create(item).Error
	}
	if err != nil {
		return err
	}

	// 已存在 → 数量累加
	return r.db.WithContext(ctx).Model(&existing).
		Update("quantity", existing.Quantity+item.Quantity).Error
}

// ListByUserID 查用户的购物车列表，同时 JOIN 商品信息
func (r *cartRepo) ListByUserID(ctx context.Context, userID uint64) ([]*cartModel.CartItem, error) {
	var items []*cartModel.CartItem
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}

	// 填充商品信息
	for _, item := range items {
		var p productModel.Product
		if err := r.db.WithContext(ctx).First(&p, item.ProductID).Error; err == nil {
			item.Product = &cartModel.ProductSnapshot{
				Name:      p.Name,
				Price:     p.Price,
				MainImage: p.MainImage,
				Stock:     p.Stock,
				Status:    p.Status,
			}
		}
	}

	return items, nil
}

func (r *cartRepo) UpdateQuantity(ctx context.Context, id, userID uint64, quantity int) error {
	return r.db.WithContext(ctx).
		Model(&cartModel.CartItem{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("quantity", quantity).Error
}

func (r *cartRepo) DeleteByID(ctx context.Context, id, userID uint64) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&cartModel.CartItem{}).Error
}

func (r *cartRepo) ClearByUserID(ctx context.Context, userID uint64) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&cartModel.CartItem{}).Error
}

func (r *cartRepo) GetByID(ctx context.Context, id, userID uint64) (*cartModel.CartItem, error) {
	var item cartModel.CartItem
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &item, err
}

func (r *cartRepo) GetByUserAndProduct(ctx context.Context, userID, productID uint64) (*cartModel.CartItem, error) {
	var item cartModel.CartItem
	err := r.db.WithContext(ctx).Where("user_id = ? AND product_id = ?", userID, productID).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &item, err
}
