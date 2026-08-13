package router

import (
	"gin-gorm-coupon-service/middleware"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// SetupRouter 只负责创建引擎 + 挂全局中间件
// 不接收任何 handler，不注册任何业务路由
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// 日志
	r.Use(middleware.LoggerMiddleware())

	// 全局限流：每秒 20，桶容量 50
	globalLimiter := middleware.NewTokenBucketLimiter(rate.Limit(20), 50)
	r.Use(globalLimiter.Middleware())

	return r
}
