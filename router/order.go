package router

import (
	"gin-gorm-coupon-service/internal/order/handler"
	"gin-gorm-coupon-service/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(r *gin.Engine, h *handler.OrderHandler) {
	orderGroup := r.Group("/api/order")
	orderGroup.Use(middleware.AuthMiddleware())
	{
		orderGroup.POST("", h.Create)
		orderGroup.GET("", h.List)
		orderGroup.GET("/:id", h.GetByID)
		orderGroup.POST("/cancel", h.Cancel)
	}
}
