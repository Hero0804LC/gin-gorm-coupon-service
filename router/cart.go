package router

import (
	"gin-gorm-coupon-service/internal/cart/handler"
	"gin-gorm-coupon-service/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterCartRoutes 注册购物车路由
func RegisterCartRoutes(r *gin.Engine, h *handler.CartHandler) {
	cartGroup := r.Group("/api/cart")
	cartGroup.Use(middleware.AuthMiddleware())
	{
		cartGroup.POST("/add", h.AddToCart)
		cartGroup.GET("/list", h.List)
		cartGroup.POST("/:id", h.UpdateQuantity)
		cartGroup.POST("/:id", h.Delete)
		cartGroup.POST("/clear", h.Clear)
	}
}
