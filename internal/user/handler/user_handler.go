package handler

import (
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
