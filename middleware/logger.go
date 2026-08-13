package middleware

import (
	"time"

	"gin-gorm-coupon-service/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		latency := time.Since(start)
		ip := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		statusCode := c.Writer.Status()
		errorMsg := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if query != "" {
			path = path + "?" + query
		}

		// logger.Logger = 你自己的全局实例
		// zap.String / zap.Int = 第三方库的函数
		log := logger.Logger.With(
			zap.String("ip", ip),
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
		)

		if statusCode >= 500 {
			log.Error("server error", zap.String("error", errorMsg))
		} else if statusCode >= 400 {
			log.Warn("client error", zap.String("error", errorMsg))
		} else {
			log.Info("request completed")
		}
	}
}
