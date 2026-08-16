package handler

import (
	"net/http"
	"strconv"

	"gin-gorm-coupon-service/internal/coupon/service"
	"gin-gorm-coupon-service/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type CouponHandler struct {
	svc *service.CouponService
}

func NewCouponHandler(svc *service.CouponService) *CouponHandler {
	return &CouponHandler{svc: svc}
}

func (h *CouponHandler) Create(c *gin.Context) {
	var req service.CreateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if err := h.svc.Create(c.Request.Context(), &req); err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, "优惠券创建成功")
}

func (h *CouponHandler) List(c *gin.Context) {
	list, err := h.svc.List(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.Success(c, list)
}

func (h *CouponHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "ID 无效")
		return
	}
	coupon, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, coupon)
}
