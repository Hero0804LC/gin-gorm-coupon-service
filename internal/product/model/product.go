package model

import "time"

type Product struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string    `gorm:"size:100;not null" json:"name"`
	Description   string    `gorm:"type:text" json:"description"`
	Price         float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	OriginalPrice float64   `gorm:"type:decimal(10,2);default:0" json:"original_price"`
	Stock         uint32    `gorm:"not null;default:0" json:"stock"`
	CategoryID    uint64    `gorm:"default:0" json:"category_id"`
	MainImage     string    `gorm:"size:255;default:''" json:"main_image"`
	Images        string    `gorm:"type:json" json:"images"`
	Status        int8      `gorm:"not null;default:1" json:"status"` // 1=上架 0=下架
	Sales         uint32    `gorm:"default:0" json:"sales"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (Product) TableName() string {
	return "product"
}
