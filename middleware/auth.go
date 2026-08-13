package middleware

import (
	"strings"

	"gin-gorm-coupon-service/internal/pkg/jwt"
	"gin-gorm-coupon-service/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT 鉴权中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//从 Header 取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Fail(c, 401, "未登录，请先登录")
			c.Abort()
			return
		}

		//格式必须是：Bearer xxxxx
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Fail(c, 401, "token 格式错误，应为 Bearer xxxx")
			c.Abort()
			return
		}

		tokenString := parts[1]

		//解析token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			response.Fail(c, 401, err.Error())
			c.Abort()
			return
		}

		//检查黑名单
		blacklisted, err := jwt.IsBlacklisted(c.Request.Context(), tokenString)
		if err != nil {
			response.Fail(c, 500, "系统错误")
			c.Abort()
			return
		}
		if blacklisted {
			response.Fail(c, 401, "token 已失效，请重新登录")
			c.Abort()
			return
		}

		//把用户信息注入 context（后续 Handler 直接用）
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
