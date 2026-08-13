package service

import (
	"context"
	"fmt"

	"gin-gorm-coupon-service/internal/product/model"
	"gin-gorm-coupon-service/internal/product/repository"
)

type ProductRepo interface {
	Create(ctx context.Context, product *model.Product) error
	GetByID(ctx context.Context, id uint64) (*model.Product, error)
	List(ctx context.Context, page, pageSize int, keyword string, categoryID uint64) ([]*model.Product, int64, error)
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id uint64) error
}

type ProductService struct {
	repo repository.ProductRepo
}

func NewProductService(repo repository.ProductRepo) *ProductService {
	return &ProductService{repo: repo}
}

// CreateProductRequest 创建商品请求
type CreateProductRequest struct {
	Name          string  `json:"name" binding:"required,min=1,max=100"`
	Description   string  `json:"description"`
	Price         float64 `json:"price" binding:"required,gt=0"`
	OriginalPrice float64 `json:"original_price"`
	Stock         uint32  `json:"stock" binding:"required"`
	CategoryID    uint64  `json:"category_id"`
	MainImage     string  `json:"main_image"`
	Images        string  `json:"images"`
}

// Create 创建商品
func (s *ProductService) Create(ctx context.Context, req *CreateProductRequest) error {
	product := &model.Product{
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		Stock:         req.Stock,
		CategoryID:    req.CategoryID,
		MainImage:     req.MainImage,
		Images:        req.Images,
		Status:        1,
	}
	return s.repo.Create(ctx, product)
}

func (s *ProductService) GetByID(ctx context.Context, id uint64) (*model.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, fmt.Errorf("商品不存在")
	}
	return product, nil
}

func (s *ProductService) List(ctx context.Context, page, pageSize int, keyword string, categoryID uint64) ([]*model.Product, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 10
	}
	return s.repo.List(ctx, page, pageSize, keyword, categoryID)
}

// UpdateProductRequest 更新商品请求
type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       uint32  `json:"stock"`
	CategoryID  uint64  `json:"category_id"`
	MainImage   string  `json:"main_image"`
	Images      string  `json:"images"`
}

// Update 更新商品
func (s *ProductService) Update(ctx context.Context, id uint64, req *UpdateProductRequest) error {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if product == nil {
		return fmt.Errorf("商品不存在")
	}

	// 只更新非空字段
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Stock > 0 {
		product.Stock = req.Stock
	}
	if req.CategoryID > 0 {
		product.CategoryID = req.CategoryID
	}
	if req.MainImage != "" {
		product.MainImage = req.MainImage
	}
	if req.Images != "" {
		product.Images = req.Images
	}

	return s.repo.Update(ctx, product)
}

func (s *ProductService) Delete(ctx context.Context, id uint64) error {
	return s.repo.Delete(ctx, id)
}
