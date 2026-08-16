package handler

import (
	"net/http"
	"strconv"

	"gin-gorm-coupon-service/internal/pkg/response"
	"gin-gorm-coupon-service/internal/seckill/service"

	"github.com/gin-gonic/gin"
)

type SeckillHandler struct {
	svc *service.SeckillService
}

func NewSeckillHandler(svc *service.SeckillService) *SeckillHandler {
	return &SeckillHandler{svc: svc}
}

func (h *SeckillHandler) Grab(c *gin.Context) {
	couponID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "优惠券 ID 无效")
		return
	}

	userID, _ := c.Get("user_id")
	// userID 从 JWT 中间件注入的是 uint64，直接类型断言
	uid, ok := userID.(uint64)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "用户未登录或 token 无效")
		return
	}

	msg, err := h.svc.Grab(c.Request.Context(), uid, couponID)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, msg)
}
