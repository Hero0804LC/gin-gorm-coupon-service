package router

import (
	"gin-gorm-coupon-service/internal/user/handler"
	"gin-gorm-coupon-service/middleware"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RegisterUserRoutes 注册用户模块路由
func RegisterUserRoutes(r *gin.Engine, h *handler.UserHandler) {

	// ========== 公开接口 ==========
	userGroup := r.Group("/api/user")
	{
		smsLimiter := middleware.NewTokenBucketLimiter(rate.Limit(1.0/60.0), 1)
		userGroup.POST("/send-code", smsLimiter.Middleware(), h.SendCode) //发送验证码
		userGroup.POST("/register", h.Register)                           //用户注册
		userGroup.POST("/login", h.Login)                                 //用户登录
	}

	// ========== 需登录接口 ==========
	authGroup := r.Group("/api/user")
	authGroup.Use(middleware.AuthMiddleware())
	{
		authGroup.GET("/profile", h.Profile) //刷新jwt
		authGroup.POST("/logout", h.Logout)  //用户登出
	}
}
