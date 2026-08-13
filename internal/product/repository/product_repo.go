package repository

import (
	"context"
	"errors"

	"gin-gorm-coupon-service/internal/product/model"

	"gorm.io/gorm"
)

type ProductRepo interface {
	Create(ctx context.Context, product *model.Product) error
	GetByID(ctx context.Context, id uint64) (*model.Product, error)
	List(ctx context.Context, page, pageSize int, keyword string, categoryID uint64) ([]*model.Product, int64, error)
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id uint64) error
}

type productRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) ProductRepo {
	return &productRepo{db: db}
}

func (r *productRepo) Create(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepo) GetByID(ctx context.Context, id uint64) (*model.Product, error) {
	var product model.Product
	err := r.db.WithContext(ctx).First(&product, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &product, err
}

func (r *productRepo) List(ctx context.Context, page, pageSize int, keyword string, categoryID uint64) ([]*model.Product, int64, error) {
	var list []*model.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Product{}).Where("status = ?", 1)

	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error

	return list, total, err
}

func (r *productRepo) Update(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&model.Product{}).Where("id = ?", id).Update("status", 0).Error
}
