package model

import "time"

type Coupon struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"size:100;not null" json:"name"`
	Type       int8      `gorm:"not null" json:"type"`
	Value      float64   `gorm:"type:decimal(10,2);not null" json:"value"`
	MinAmount  float64   `gorm:"type:decimal(10,2)" json:"min_amount"`
	TotalCount int       `gorm:"not null" json:"total_count"`
	UsedCount  int       `gorm:"not null;default:0" json:"used_count"`
	PerLimit   int       `gorm:"not null" json:"per_limit"`
	StartTime  time.Time `gorm:"type:datetime" json:"start_time"`
	EndTime    time.Time `gorm:"type:datetime" json:"end_time"`
	Status     int8      `gorm:"not null;default:1" json:"status"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Coupon) TableName() string { return "coupon" }

type UserCoupon struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64     `gorm:"not null;index:idx_user_coupon,unique" json:"user_id"`
	CouponID  uint64     `gorm:"not null;index:idx_user_coupon,unique" json:"coupon_id"`
	OrderID   uint64     `gorm:"default:0" json:"order_id"`
	Status    int8       `gorm:"not null;default:0" json:"status"`
	UsedAt    *time.Time `json:"used_at"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (UserCoupon) TableName() string { return "user_coupon" }
