package service

import (
	"context"
	"fmt"

	cartModel "gin-gorm-coupon-service/internal/cart/model"
	cartRepo "gin-gorm-coupon-service/internal/cart/repository"
	productModel "gin-gorm-coupon-service/internal/product/model"
	productRepo "gin-gorm-coupon-service/internal/product/repository"
)

type CartRepo interface {
	AddOrUpdate(ctx context.Context, item *cartModel.CartItem) error
	ListByUserID(ctx context.Context, userID uint64) ([]*cartModel.CartItem, error)
	UpdateQuantity(ctx context.Context, id, userID uint64, quantity int) error
	DeleteByID(ctx context.Context, id, userID uint64) error
	ClearByUserID(ctx context.Context, userID uint64) error
}

type ProductRepo interface {
	GetByID(ctx context.Context, id uint64) (*productModel.Product, error)
}

type CartService struct {
	cartRepo    cartRepo.CartRepo
	productRepo productRepo.ProductRepo
}

func NewCartService(cartRepo cartRepo.CartRepo, productRepo productRepo.ProductRepo) *CartService {
	return &CartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

type AddToCartRequest struct {
	ProductID uint64 `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

// AddToCart 加入购物车
func (s *CartService) AddToCart(ctx context.Context, userID uint64, req *AddToCartRequest) error {
	// 1. 校验商品是否存在且上架
	p, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return fmt.Errorf("查询商品失败")
	}
	if p == nil || p.Status != 1 {
		return fmt.Errorf("商品不存在或已下架")
	}

	// 2. 校验数量
	if req.Quantity <= 0 {
		return fmt.Errorf("数量必须大于 0")
	}
	if uint32(req.Quantity) > p.Stock {
		return fmt.Errorf("库存不足")
	}

	// 3. 加入购物车（已存在则数量累加）
	item := &cartModel.CartItem{
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Checked:   1,
	}
	return s.cartRepo.AddOrUpdate(ctx, item)
}

// List 查看购物车
func (s *CartService) List(ctx context.Context, userID uint64) ([]*cartModel.CartItem, error) {
	return s.cartRepo.ListByUserID(ctx, userID)
}

// UpdateQuantity 修改数量
func (s *CartService) UpdateQuantity(ctx context.Context, userID, id uint64, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("数量必须大于 0")
	}

	// 校验购物车项是否属于当前用户
	item, err := s.cartRepo.GetByID(ctx, id, userID)
	if err != nil {
		return fmt.Errorf("查询失败")
	}
	if item == nil {
		return fmt.Errorf("购物车项不存在")
	}

	// 校验库存
	p, _ := s.productRepo.GetByID(ctx, item.ProductID)
	if p != nil && uint32(quantity) > p.Stock {
		return fmt.Errorf("库存不足")
	}

	return s.cartRepo.UpdateQuantity(ctx, id, userID, quantity)
}

// Delete 删除单项
func (s *CartService) Delete(ctx context.Context, userID, id uint64) error {
	return s.cartRepo.DeleteByID(ctx, id, userID)
}

// Clear 清空购物车
func (s *CartService) Clear(ctx context.Context, userID uint64) error {
	return s.cartRepo.ClearByUserID(ctx, userID)
}
