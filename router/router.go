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
		userGroup.POST("/send-code", userHandler.SendCode) //发送验证码
		userGroup.POST("/register", userHandler.Register)  //用户注册
	}

	return r
}
