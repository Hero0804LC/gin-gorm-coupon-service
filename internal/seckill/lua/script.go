package lua

import (
	"context"
	_ "embed"

	"github.com/redis/go-redis/v9"
)

//go:embed grab_coupon.lua
var grabCouponLua string

var GrabCouponScript = redis.NewScript(grabCouponLua)

// Load 预加载脚本到 Redis（可选，提高性能）
func Load(ctx context.Context, rdb *redis.Client) error {
	return GrabCouponScript.Load(ctx, rdb).Err()
}
