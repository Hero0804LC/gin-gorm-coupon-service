package router

import (
	"gin-gorm-coupon-service/internal/user/handler"
	"gin-gorm-coupon-service/middleware"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func SetupRouter(userHandler *handler.UserHandler) *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.Use(middleware.LoggerMiddleware())
	globalLimiter := middleware.NewTokenBucketLimiter(
		rate.Limit(20),
		50,
	)
	r.Use(globalLimiter.Middleware())

	userGroup := r.Group("/api/user")
	{
		// 发短信：单独限流（1次/分钟/IP）
		smsLimiter := middleware.NewTokenBucketLimiter(
			rate.Limit(1.0/60.0), // 每分钟 1 个
			1,                    // 桶容量 1（不允许突发）
		)
		userGroup.POST("/send-code", smsLimiter.Middleware(), userHandler.SendCode) //发送验证码
		userGroup.POST("/register", userHandler.Register)                           //用户注册
		userGroup.POST("/login", userHandler.Login)                                 //用户登录
	}
	authGroup := r.Group("/api/user")
	authGroup.Use(middleware.AuthMiddleware())
	{
		authGroup.GET("/profile", userHandler.Profile) //获取用户信息
		authGroup.POST("/logout", userHandler.Logout)  //退出登录
	}

	return r
}
