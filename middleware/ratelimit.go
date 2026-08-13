package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"gin-gorm-coupon-service/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// TokenBucketLimiter 基于 IP 的令牌桶限流器
type TokenBucketLimiter struct {
	limiters sync.Map // key: IP, value: *rate.Limiter
	r        rate.Limit
	burst    int
}

// NewTokenBucketLimiter 创建限流器
func NewTokenBucketLimiter(r rate.Limit, burst int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		r:     r,     //每秒填充速率
		burst: burst, //桶容量
	}
}

// getLimiter 获取或创建 IP 对应的限流器
func (t *TokenBucketLimiter) getLimiter(ip string) *rate.Limiter {
	if v, ok := t.limiters.Load(ip); ok {
		return v.(*rate.Limiter)
	}

	limiter := rate.NewLimiter(t.r, t.burst)
	t.limiters.Store(ip, limiter)
	return limiter
}

// Cleanup 定期清理不活跃的限流器（防止内存泄漏）
func (t *TokenBucketLimiter) Cleanup(ctx context.Context, interval, expire time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				time.Sleep(expire)
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Middleware 返回 Gin 中间件
func (t *TokenBucketLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := t.getLimiter(ip)

		if !limiter.Allow() {
			response.Fail(c, http.StatusTooManyRequests, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}
