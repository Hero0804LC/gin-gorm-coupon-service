package model

import "time"

// Order 订单主表
type Order struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNo      string     `gorm:"size:64;uniqueIndex;not null" json:"order_no"`
	UserID       uint64     `gorm:"not null;index" json:"user_id"`
	TotalAmount  float64    `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	PayAmount    float64    `gorm:"type:decimal(10,2);not null" json:"pay_amount"`
	CouponID     uint64     `gorm:"default:0" json:"coupon_id"`
	CouponAmount float64    `gorm:"type:decimal(10,2);default:0" json:"coupon_amount"`
	Status       int8       `gorm:"not null;default:0" json:"status"` // 0=待支付 1=已支付 2=已发货 3=已完成 4=已取消 5=已退款
	AddressID    uint64     `gorm:"default:0" json:"address_id"`
	Remark       string     `gorm:"size:200" json:"remark"`
	PaidAt       *time.Time `json:"paid_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	Items []*OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

func (Order) TableName() string { return "order" }

// OrderItem 订单明细表
type OrderItem struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID      uint64    `gorm:"not null;index" json:"order_id"`
	OrderNo      string    `gorm:"size:64;index" json:"order_no"`
	ProductID    uint64    `gorm:"not null" json:"product_id"`
	ProductName  string    `gorm:"size:100" json:"product_name"`
	ProductImage string    `gorm:"size:255" json:"product_image"`
	Price        float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Quantity     int       `gorm:"not null" json:"quantity"`
	TotalPrice   float64   `gorm:"type:decimal(10,2);not null" json:"total_price"`
	CreatedAt    time.Time `json:"created_at"`
}

func (OrderItem) TableName() string { return "order_item" }
