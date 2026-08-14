package handler

import (
	"net/http"
	"strconv"

	"gin-gorm-coupon-service/internal/cart/service"
	"gin-gorm-coupon-service/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	svc *service.CartService
}

func NewCartHandler(svc *service.CartService) *CartHandler {
	return &CartHandler{svc: svc}
}

// AddToCart 加入购物车
func (h *CartHandler) AddToCart(c *gin.Context) {
	var req service.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	userID, _ := c.Get("user_id")
	if err := h.svc.AddToCart(c.Request.Context(), userID.(uint64), &req); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, "加入购物车成功")
}

// List 购物车列表
func (h *CartHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	list, err := h.svc.List(c.Request.Context(), userID.(uint64))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}

	response.Success(c, list)
}

// UpdateQuantity 修改数量
func (h *CartHandler) UpdateQuantity(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "ID 无效")
		return
	}

	var req struct {
		Quantity int `json:"quantity" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}

	userID, _ := c.Get("user_id")
	if err := h.svc.UpdateQuantity(c.Request.Context(), userID.(uint64), id, req.Quantity); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, "修改成功")
}

// Delete 删除单项
func (h *CartHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "ID 无效")
		return
	}

	userID, _ := c.Get("user_id")
	if err := h.svc.Delete(c.Request.Context(), userID.(uint64), id); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, "删除成功")
}

// Clear 清空购物车
func (h *CartHandler) Clear(c *gin.Context) {
	userID, _ := c.Get("user_id")
	if err := h.svc.Clear(c.Request.Context(), userID.(uint64)); err != nil {
		response.Fail(c, http.StatusInternalServerError, "清空失败")
		return
	}

	response.Success(c, "清空成功")
}
