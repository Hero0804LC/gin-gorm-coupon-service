package handler

import (
	"gin-gorm-coupon-service/internal/pkg/jwt"
	"gin-gorm-coupon-service/internal/user/service"

	"gin-gorm-coupon-service/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// SendCodeRequest 发送验证码请求
type SendCodeRequest struct {
	Phone string `json:"phone" binding:"required,len=11"`
}

// SendCode 发送验证码接口
func (h *UserHandler) SendCode(c *gin.Context) {
	var req SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := h.userService.SendCode(c.Request.Context(), req.Phone); err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, "验证码已发送")
}

// RegisterRequest 注册用户请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone" binding:"required,len=11"`
	Code     string `json:"code" binding:"required,len=6"`
}

// Register 注册用户
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := h.userService.Register(
		c.Request.Context(),
		req.Username,
		req.Password,
		req.Phone,
		req.Code,
	); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, "注册成功")
}

type LoginRequest struct {
	Phone    string `json:"phone" binding:"required,len=11"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	token, err := h.userService.Login(
		c.Request.Context(),
		req.Phone,
		req.Password,
	)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, LoginResponse{Token: token})
}

// Profile 获取当前用户信息（需要登录）
func (h *UserHandler) Profile(c *gin.Context) {
	// 从中间件注入的 context 取 user_id
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Fail(c, 401, "未登录")
		return
	}

	// 类型断言（JWT 里存的是 uint64）
	userID, ok := userIDVal.(uint64)
	if !ok {
		response.Fail(c, 401, "用户信息异常")
		return
	}

	user, err := h.userService.Profile(c.Request.Context(), userID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 不返回密码
	user.Password = ""

	response.Success(c, user)
}

type LogoutRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
}

// Logout 退出登录
func (h *UserHandler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	// token 拿到过期时间
	claims, err := jwt.ParseToken(req.AccessToken)
	if err != nil {
		// token 已经无效/过期，直接算退出成功
		response.Success(c, "退出成功")
		return
	}

	//加入黑名单
	if err := jwt.AddToBlacklist(
		c.Request.Context(),
		req.AccessToken,
		claims.ExpiresAt.Time,
	); err != nil {
		response.Fail(c, http.StatusInternalServerError, "退出失败")
		return
	}

	response.Success(c, "退出成功")
}
