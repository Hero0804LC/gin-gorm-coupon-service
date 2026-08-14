package handler

import (
	"net/http"
	"strconv"

	"gin-gorm-coupon-service/internal/order/service"
	"gin-gorm-coupon-service/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	svc *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	userID, _ := c.Get("user_id")
	order, err := h.svc.CreateFromCart(c.Request.Context(), userID.(uint64), &req)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{
		"order_no":   order.OrderNo,
		"pay_amount": order.PayAmount,
		"status":     order.Status,
	})
}

func (h *OrderHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	userID, _ := c.Get("user_id")
	list, total, err := h.svc.List(c.Request.Context(), userID.(uint64), page, pageSize)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}

	response.Success(c, gin.H{
		"list":  list,
		"total": total,
		"page":  page,
		"size":  pageSize,
	})
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "ID 无效")
		return
	}

	userID, _ := c.Get("user_id")
	order, err := h.svc.GetByID(c.Request.Context(), id, userID.(uint64))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, order)
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	var req struct {
		OrderNo string `json:"order_no" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	userID, _ := c.Get("user_id")
	if err := h.svc.Cancel(c.Request.Context(), req.OrderNo, userID.(uint64)); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, "取消成功")
}
