package main

import (
	"gin-gorm-coupon-service/app"
	"gin-gorm-coupon-service/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 1. 配置
	if err := config.Init("config/config.yaml"); err != nil {
		log.Fatalf("config init failed: %v", err)
	}
	cfg := config.Get()

	// 2. MySQL
	db, err := gorm.Open(mysql.Open(cfg.MySQL.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("mysql init failed: %v", err)
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// 3. 初始化应用（✅ 依赖注入已剥离）
	application, err := app.NewApp(cfg, db)
	if err != nil {
		log.Fatalf("app init failed: %v", err)
	}
	defer application.Close()

	// 4. 启动
	go func() {
		log.Println("server listening on", cfg.Server.Addr)
		if err := application.Router.Run(cfg.Server.Addr); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	// 5. 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("server exiting")
}
