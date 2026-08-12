package app

import (
	"gin-gorm-coupon-service/config"
	"gin-gorm-coupon-service/internal/pkg/redis"
	"gin-gorm-coupon-service/internal/user/handler"
	"gin-gorm-coupon-service/internal/user/repository"
	"gin-gorm-coupon-service/internal/user/service"
	"gin-gorm-coupon-service/router"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// App 应用上下文
type App struct {
	Router *gin.Engine
	DB     *gorm.DB
}

// NewApp 初始化所有依赖
func NewApp(cfg *config.Config, db *gorm.DB) (*App, error) {
	// 1. 基础设施
	if err := redis.Init(); err != nil {
		return nil, err
	}

	// 2. User 模块
	userRepo := repository.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// 3. 路由
	r := router.SetupRouter(userHandler)

	return &App{
		Router: r,
		DB:     db,
	}, nil
}

// Close 释放资源
func (a *App) Close() {
	redis.Close()
}
