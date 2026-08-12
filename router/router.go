package router

import (
	"gin-gorm-coupon-service/internal/user/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *handler.UserHandler) *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	userGroup := r.Group("/api/user")
	{
		userGroup.POST("/send-code", userHandler.SendCode)
	}

	return r
}
