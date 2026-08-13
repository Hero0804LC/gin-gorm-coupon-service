package jwt

import (
	"context"
	"errors"
	"gin-gorm-coupon-service/internal/pkg/redis"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	ErrTokenExpired = errors.New("token 已过期")
	ErrTokenInvalid = errors.New("token 无效")
)

// Claims 自定义载荷
type Claims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Config JWT 配置
type Config struct {
	Secret string
	Expire time.Duration
}

var cfg Config

// 黑名单
const blacklistPrefix = "jwt:blacklist:"

// Init 初始化
func Init(secret string, expire time.Duration) {
	cfg = Config{
		Secret: secret,
		Expire: expire,
	}
}

// GenerateToken 生成 Token
func GenerateToken(userID uint64, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.Expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

// ParseToken 解析 Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// AddToBlacklist 将 token 加入黑名单
// tokenString: 原始 token 字符串
// exp: token 的过期时间（从 claims 里取）
func AddToBlacklist(ctx context.Context, tokenString string, exp time.Time) error {
	key := blacklistPrefix + tokenString

	// TTL = token 剩余有效期
	ttl := time.Until(exp)
	if ttl <= 0 {
		return nil // 已经过期，不用加
	}

	return redis.Client.Set(ctx, key, "1", ttl).Err()
}

// IsBlacklisted 检查 token 是否在黑名单中
func IsBlacklisted(ctx context.Context, tokenString string) (bool, error) {
	key := blacklistPrefix + tokenString
	exists, err := redis.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}
