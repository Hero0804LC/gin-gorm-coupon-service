-- KEYS[1] = 库存 key: seckill:stock:{coupon_id}
-- KEYS[2] = 用户限领 key: seckill:user_limit:{coupon_id}:{user_id}
-- ARGV[1] = 总库存上限
-- ARGV[2] = 每人限领数量
-- ARGV[3] = 当前时间戳（用于判断是否过期）

local stock_key = KEYS[1]
local user_key = KEYS[2]
local total = tonumber(ARGV[1])
local per_limit = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

-- 1. 检查库存
local stock = redis.call('GET', stock_key)
if not stock then
    -- 首次初始化库存
    redis.call('SET', stock_key, total)
    stock = total
end

if tonumber(stock) <= 0 then
    return -1  -- 库存不足
end

-- 2. 检查用户限领
local user_count = redis.call('GET', user_key)
if user_count and tonumber(user_count) >= per_limit then
    return -2  -- 超过限领数量
end

-- 3. 原子扣减
redis.call('DECR', stock_key)
redis.call('INCR', user_key)
-- 设置用户限领 key 过期时间（24小时，防止永久占用内存）
redis.call('EXPIRE', user_key, 86400)

return 1  -- 成功