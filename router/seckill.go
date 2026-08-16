package router

import (
	"gin-gorm-coupon-service/internal/seckill/handler"
	"gin-gorm-coupon-service/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterSeckillRoutes(r *gin.Engine, h *handler.SeckillHandler) {
	sg := r.Group("/api/seckill")
	sg.Use(middleware.AuthMiddleware())
	{
		sg.POST("/grab/:id", h.Grab)
	}
}
