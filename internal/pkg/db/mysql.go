package db

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMySQL(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("mysql init failed: %v", err)
	}
	return db
}
