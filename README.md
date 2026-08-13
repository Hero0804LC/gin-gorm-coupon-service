# gin-gorm-coupon-service

### 用户模块
- 注册：接收手机号与密码，Redis限制60s内不可重复发送验证码，密码经bcrypt哈希后写入 MySQL
- 登录：校验账号密码，成功后签发 JWT Token

### 中间件
- 鉴权：请求携带Token经Gin中间件解析并校验有效性，通过后将 user_id 注入上下文
- 限流：基于令牌桶算法实现接口级限流，超出阈值返回 429 Too Many Requests