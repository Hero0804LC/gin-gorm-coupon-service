package router

import (
	"gin-gorm-coupon-service/internal/user/handler"
	"gin-gorm-coupon-service/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *handler.UserHandler) *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	userGroup := r.Group("/api/user")
	{
		userGroup.POST("/send-code", userHandler.SendCode) //发送验证码
		userGroup.POST("/register", userHandler.Register)  //用户注册
		userGroup.POST("/login", userHandler.Login)        //用户登录
	}
	authGroup := r.Group("/api/user")
	authGroup.Use(middleware.AuthMiddleware())
	{
		authGroup.GET("/profile", userHandler.Profile) //获取用户信息
	}

	return r
}
