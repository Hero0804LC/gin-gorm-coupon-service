package model

import "time"

type CartItem struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64    `gorm:"not null;index:idx_user_id" json:"user_id"`
	ProductID uint64    `gorm:"not null" json:"product_id"`
	Quantity  int       `gorm:"not null;default:1" json:"quantity"`
	Checked   int8      `gorm:"not null;default:1" json:"checked"` // 1=选中 0=未选中
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Product *ProductSnapshot `gorm:"-" json:"product,omitempty"`
}

// ProductSnapshot 购物车列表里展示用的商品快照
type ProductSnapshot struct {
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	MainImage string  `json:"main_image"`
	Stock     uint32  `json:"stock"`
	Status    int8    `json:"status"`
}

func (CartItem) TableName() string {
	return "cart_item"
}
