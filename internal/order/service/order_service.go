package service

import (
	"context"
	"fmt"
	cartModel "gin-gorm-coupon-service/internal/cart/model"
	orderModel "gin-gorm-coupon-service/internal/order/model"
	"gin-gorm-coupon-service/internal/order/repository"
	"gin-gorm-coupon-service/internal/order/utils"
	productModel "gin-gorm-coupon-service/internal/product/model"
	productRepo "gin-gorm-coupon-service/internal/product/repository"
)

type OrderRepo interface {
	CreateOrder(ctx context.Context, order *orderModel.Order, items []*orderModel.OrderItem) error
	GetByOrderNo(ctx context.Context, orderNo string) (*orderModel.Order, error)
	ListByUserID(ctx context.Context, userID uint64, page, pageSize int) ([]*orderModel.Order, int64, error)
	GetByID(ctx context.Context, id, userID uint64) (*orderModel.Order, error)
	UpdateStatus(ctx context.Context, orderNo string, status int8) error
}

type ProductRepo interface {
	GetByID(ctx context.Context, id uint64) (*productModel.Product, error)
}

type CartRepo interface {
	ListByUserID(ctx context.Context, userID uint64) ([]*cartModel.CartItem, error)
	ClearByUserID(ctx context.Context, userID uint64) error
}

type OrderService struct {
	orderRepo   repository.OrderRepo
	productRepo productRepo.ProductRepo
	cartRepo    CartRepo
}

func NewOrderService(orderRepo repository.OrderRepo, productRepo productRepo.ProductRepo, cartRepo CartRepo) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		cartRepo:    cartRepo,
	}
}

type CreateOrderRequest struct {
	AddressID    uint64  `json:"address_id"`
	CouponID     uint64  `json:"coupon_id"`
	CouponAmount float64 `json:"coupon_amount"`
	Remark       string  `json:"remark"`
}

// CreateFromCart 从购物车创建订单
func (s *OrderService) CreateFromCart(ctx context.Context, userID uint64, req *CreateOrderRequest) (*orderModel.Order, error) {
	cartItems, err := s.cartRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("获取购物车失败")
	}
	// 筛选选中的
	var selectedItems []*cartModel.CartItem
	for _, item := range cartItems {
		if item.Checked == 1 {
			selectedItems = append(selectedItems, item)
		}
	}
	if len(selectedItems) == 0 {
		return nil, fmt.Errorf("购物车为空，请先添加商品")
	}

	//构建订单明细
	var orderItems []*orderModel.OrderItem
	var totalAmount float64

	for _, ci := range selectedItems {
		p, err := s.productRepo.GetByID(ctx, ci.ProductID)
		if err != nil || p == nil || p.Status != 1 {
			return nil, fmt.Errorf("商品 [%d] 不存在或已下架", ci.ProductID)
		}
		if uint32(ci.Quantity) > p.Stock {
			return nil, fmt.Errorf("商品 [%s] 库存不足", p.Name)
		}

		itemTotal := p.Price * float64(ci.Quantity)
		totalAmount += itemTotal

		orderItems = append(orderItems, &orderModel.OrderItem{
			ProductID:    ci.ProductID,
			ProductName:  p.Name,
			ProductImage: p.MainImage,
			Price:        p.Price,
			Quantity:     ci.Quantity,
			TotalPrice:   itemTotal,
		})
	}
	//计算实付金额
	payAmount := totalAmount - req.CouponAmount
	if payAmount < 0 {
		payAmount = 0
	}
	//生成订单号
	orderNo := utils.GenerateOrderNo()
	//创建订单
	order := &orderModel.Order{
		OrderNo:      orderNo,
		UserID:       userID,
		TotalAmount:  totalAmount,
		PayAmount:    payAmount,
		CouponID:     req.CouponID,
		CouponAmount: req.CouponAmount,
		Status:       0, // 待支付
		AddressID:    req.AddressID,
		Remark:       req.Remark,
	}
	//事务写入（订单 + 明细 + 扣库存）
	if err := s.orderRepo.CreateOrder(ctx, order, orderItems); err != nil {
		return nil, fmt.Errorf("创建订单失败: %w", err)
	}

	//清空购物车（已下单的商品）
	_ = s.cartRepo.ClearByUserID(ctx, userID)

	return order, nil
}

// List 订单列表
func (s *OrderService) List(ctx context.Context, userID uint64, page, pageSize int) ([]*orderModel.Order, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 10
	}
	return s.orderRepo.ListByUserID(ctx, userID, page, pageSize)
}

// GetByID 订单详情
func (s *OrderService) GetByID(ctx context.Context, id, userID uint64) (*orderModel.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, fmt.Errorf("订单不存在")
	}
	return order, nil
}

// Cancel 取消订单
func (s *OrderService) Cancel(ctx context.Context, orderNo string, userID uint64) error {
	order, err := s.orderRepo.GetByOrderNo(ctx, orderNo)
	if err != nil || order == nil {
		return fmt.Errorf("订单不存在")
	}
	if order.UserID != userID {
		return fmt.Errorf("无权限操作")
	}
	if order.Status != 0 {
		return fmt.Errorf("当前订单状态不支持取消")
	}
	return s.orderRepo.UpdateStatus(ctx, orderNo, 4) // 4=已取消
}
