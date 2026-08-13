package app

import (
	"fmt"
	"gin-gorm-coupon-service/config"
	"gin-gorm-coupon-service/internal/pkg/jwt"
	"github.com/gin-gonic/gin"

	"gin-gorm-coupon-service/internal/pkg/logger"
	"gin-gorm-coupon-service/internal/pkg/redis"
	// 用户模块
	userHandler "gin-gorm-coupon-service/internal/user/handler"
	userRepo "gin-gorm-coupon-service/internal/user/repository"
	userService "gin-gorm-coupon-service/internal/user/service"

	// 商品模块
	productHandler "gin-gorm-coupon-service/internal/product/handler"
	productRepo "gin-gorm-coupon-service/internal/product/repository"
	productService "gin-gorm-coupon-service/internal/product/service"
	"gin-gorm-coupon-service/router"

	"gorm.io/gorm"
)

type App struct {
	Router *gin.Engine
	DB     *gorm.DB
}

func NewApp(db *gorm.DB) (*App, error) {

	//初始化
	logger.Init()
	if err := redis.Init(); err != nil {
		return nil, fmt.Errorf("redis init failed: %w", err)
	}

	jwt.Init(config.Get().JWT.Secret, config.Get().JWT.Expire)

	//User 模块
	userRepo := userRepo.NewUserRepo(db)
	userSvc := userService.NewUserService(userRepo)
	userHandler := userHandler.NewUserHandler(userSvc)

	//Product 模块
	productRepo := productRepo.NewProductRepo(db)
	productSvc := productService.NewProductService(productRepo)
	productHandler := productHandler.NewProductHandler(productSvc)

	// Router（先创建引擎 → 再逐个注册模块）
	r := router.SetupRouter()
	router.RegisterUserRoutes(r, userHandler)
	router.RegisterProductRoutes(r, productHandler)

	return &App{
		Router: r,
		DB:     db,
	}, nil
}
