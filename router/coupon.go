package router

import (
	"gin-gorm-coupon-service/internal/coupon/handler"
	"gin-gorm-coupon-service/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterCouponRoutes(r *gin.Engine, h *handler.CouponHandler) {
	couponGroup := r.Group("/api/coupon")
	{
		couponGroup.GET("", h.List)
		couponGroup.GET("/:id", h.GetByID)
	}

	authGroup := r.Group("/api/coupon")
	authGroup.Use(middleware.AuthMiddleware())
	{
		authGroup.POST("", h.Create)
	}
}
