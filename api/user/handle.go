package user

import (
	"gin-gorm-coupon-service/internal/pkg/response"
	"gin-gorm-coupon-service/internal/user/cache"
	"gin-gorm-coupon-service/internal/user/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func SendSmsHandler(c *gin.Context) {
	var req SendSmsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	code := service.GenerateCode()

	// 存入 Redis，5分钟过期
	err := cache.SetCode(c, req.Phone, code, 5*time.Minute)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "发送失败")
		return
	}

	// 模拟发送短信
	log.Printf("模拟发送短信：phone=%s, code=%s\n", req.Phone, code)

	response.Success(c, gin.H{"message": "验证码已发送"})
}
func RegisterHandler(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 1. 校验验证码
	cacheCode, err := cache.GetCode(c, req.Phone)
	if err != nil || cacheCode != req.Code {
		response.Fail(c, http.StatusBadRequest, "验证码错误")
		return
	}

	// 2. 校验用户名/手机号是否已存在（调用 repo）
	exists, _ := repository.UserExists(req.Username, req.Phone)
	if exists {
		response.Fail(c, http.StatusBadRequest, "用户名或手机号已存在")
		return
	}

	// 3. 密码加密
	hashPwd, _ := password.Hash(req.Password)

	// 4. 写库
	user := &repository.User{
		Username: req.Username,
		Password: hashPwd,
		Phone:    req.Phone,
	}
	err = repository.CreateUser(user)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "注册失败")
		return
	}

	// 5. 删除验证码（一次性）
	_ = cache.DelCode(c, req.Phone)

	response.Success(c, gin.H{"message": "注册成功"})
}
