package router

import (
	"gin-gorm-coupon-service/internal/product/handler"
	"gin-gorm-coupon-service/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterProductRoutes 注册商品模块路由
func RegisterProductRoutes(r *gin.Engine, h *handler.ProductHandler) {

	// ========== 公开接口 ==========
	productGroup := r.Group("/api/product")
	{
		productGroup.GET("", h.List)
		productGroup.GET("/:id", h.GetByID)
	}

	// ========== 需登录接口 ==========
	authGroup := r.Group("/api/product")
	authGroup.Use(middleware.AuthMiddleware())
	{
		authGroup.POST("", h.Create)
		authGroup.PUT("/:id", h.Update)
		authGroup.DELETE("/:id", h.Delete)
	}
}
