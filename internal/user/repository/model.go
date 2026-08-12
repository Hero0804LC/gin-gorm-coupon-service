package repository

import "time"

type User struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement;column:id"`
	Username  string     `gorm:"type:varchar(50);not null;uniqueIndex:uk_username;column:username"`
	Password  string     `gorm:"type:varchar(255);not null;column:password"`
	Phone     string     `gorm:"type:varchar(20);not null;uniqueIndex:uk_phone;column:phone"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index"`
}

func (User) TableName() string {
	return "user"
}
