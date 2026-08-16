package app

import (
	"context"
	"fmt"

	"gin-gorm-coupon-service/config"
	"gin-gorm-coupon-service/internal/pkg/jwt"
	"gin-gorm-coupon-service/internal/pkg/logger"
	"gin-gorm-coupon-service/internal/pkg/redis"
	"gin-gorm-coupon-service/internal/seckill/async"
	"gin-gorm-coupon-service/router"

	// ========== User 模块 ==========
	userHandler "gin-gorm-coupon-service/internal/user/handler"
	userRepo "gin-gorm-coupon-service/internal/user/repository"
	userSvc "gin-gorm-coupon-service/internal/user/service"

	// ========== Product 模块 ==========
	productHandler "gin-gorm-coupon-service/internal/product/handler"
	productRepo "gin-gorm-coupon-service/internal/product/repository"
	productSvc "gin-gorm-coupon-service/internal/product/service"

	// ========== Cart 模块 ==========
	cartHandler "gin-gorm-coupon-service/internal/cart/handler"
	cartRepo "gin-gorm-coupon-service/internal/cart/repository"
	cartSvc "gin-gorm-coupon-service/internal/cart/service"

	// ========== Order 模块 ==========
	orderHandler "gin-gorm-coupon-service/internal/order/handler"
	orderRepo "gin-gorm-coupon-service/internal/order/repository"
	orderSvc "gin-gorm-coupon-service/internal/order/service"

	// ========== Coupon 模块 ==========
	couponHandler "gin-gorm-coupon-service/internal/coupon/handler"
	couponRepo "gin-gorm-coupon-service/internal/coupon/repository"
	couponSvc "gin-gorm-coupon-service/internal/coupon/service"

	// ========== Seckill 模块 ==========
	seckillHandler "gin-gorm-coupon-service/internal/seckill/handler"
	seckillSvc "gin-gorm-coupon-service/internal/seckill/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	Router *gin.Engine
	DB     *gorm.DB
}

func NewApp(db *gorm.DB) (*App, error) {
	// ========== 全局初始化 ==========
	logger.Init()
	if err := redis.Init(); err != nil {
		return nil, fmt.Errorf("redis init failed: %w", err)
	}
	jwt.Init(config.Get().JWT.Secret, config.Get().JWT.Expire)

	// ========== User ==========
	userRepoInst := userRepo.NewUserRepo(db)
	userSvcInst := userSvc.NewUserService(userRepoInst)
	userHandlerInst := userHandler.NewUserHandler(userSvcInst)

	// ========== Product ==========
	productRepoInst := productRepo.NewProductRepo(db)
	productSvcInst := productSvc.NewProductService(productRepoInst)
	productHandlerInst := productHandler.NewProductHandler(productSvcInst)

	// ========== Cart ==========
	cartRepoInst := cartRepo.NewCartRepo(db)
	cartSvcInst := cartSvc.NewCartService(cartRepoInst, productRepoInst)
	cartHandlerInst := cartHandler.NewCartHandler(cartSvcInst)

	// ========== Order ==========
	orderRepoInst := orderRepo.NewOrderRepo(db)
	orderSvcInst := orderSvc.NewOrderService(orderRepoInst, productRepoInst, cartRepoInst)
	orderHandlerInst := orderHandler.NewOrderHandler(orderSvcInst)

	// ========== Coupon ==========
	couponRepoInst := couponRepo.NewCouponRepo(db)
	userCouponRepoInst := couponRepo.NewUserCouponRepo(db)
	couponSvcInst := couponSvc.NewCouponService(couponRepoInst, userCouponRepoInst)
	couponHandlerInst := couponHandler.NewCouponHandler(couponSvcInst)

	// ========== Seckill ==========
	seckillSvcInst := seckillSvc.NewSeckillService(
		db,
		couponRepoInst,
		orderRepoInst,
		userCouponRepoInst,
	)

	seckillHandlerInst := seckillHandler.NewSeckillHandler(seckillSvcInst)

	worker := async.NewWorker(async.TaskChan, seckillSvcInst)
	worker.Start(context.Background())

	// ========== Router ==========
	r := router.SetupRouter()
	router.RegisterUserRoutes(r, userHandlerInst)
	router.RegisterProductRoutes(r, productHandlerInst)
	router.RegisterCartRoutes(r, cartHandlerInst)
	router.RegisterOrderRoutes(r, orderHandlerInst)
	router.RegisterCouponRoutes(r, couponHandlerInst)
	router.RegisterSeckillRoutes(r, seckillHandlerInst)

	// Worker 启动（不阻塞 HTTP）
	worker.Start(context.Background())

	return &App{
		Router: r,
		DB:     db,
	}, nil
}
