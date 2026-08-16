package service

import (
	"context"
	"fmt"
	"time"

	"gin-gorm-coupon-service/internal/coupon/model"
	orderModel "gin-gorm-coupon-service/internal/order/model"
	"gin-gorm-coupon-service/internal/seckill/async"
	"gorm.io/gorm"
)

// ========== 接口定义 ==========

type CouponRepo interface {
	GetByID(ctx context.Context, id int64) (*model.Coupon, error)
}

type OrderRepo interface {
	Create(ctx context.Context, order *orderModel.Order) error
}

type UserCouponRepo interface {
	Create(ctx context.Context, uc *model.UserCoupon) error
}

// ========== Service ==========

type SeckillService struct {
	db             *gorm.DB
	couponRepo     CouponRepo
	orderRepo      OrderRepo
	userCouponRepo UserCouponRepo
}

func NewSeckillService(
	db *gorm.DB,
	couponRepo CouponRepo,
	orderRepo OrderRepo,
	userCouponRepo UserCouponRepo,
) *SeckillService {
	return &SeckillService{
		db:             db,
		couponRepo:     couponRepo,
		orderRepo:      orderRepo,
		userCouponRepo: userCouponRepo,
	}
}

// Grab HTTP 入口
func (s *SeckillService) Grab(ctx context.Context, userID uint64, couponID uint64) (string, error) {
	coupon, err := s.couponRepo.GetByID(ctx, int64(couponID))
	if err != nil {
		return "", fmt.Errorf("优惠券不存在")
	}

	_ = coupon.StartTime
	_ = coupon.EndTime

	ok := async.Dispatch(&async.SeckillTask{
		UserID:   userID,
		CouponID: couponID,
	})
	if !ok {
		return "", fmt.Errorf("秒杀系统繁忙，请稍后再试")
	}

	return "抢购请求已接收，正在处理中", nil
}

// ProcessTask Worker 调用的处理函数（签名匹配 async.TaskHandler）
func (s *SeckillService) ProcessTask(ctx context.Context, task *async.SeckillTask) error {
	return s.db.Transaction(func(tx *gorm.DB) error {

		// 1. 查优惠券快照（事务内重新查，保证一致性）
		coupon, err := s.couponRepo.GetByID(ctx, int64(task.CouponID))
		if err != nil {
			return fmt.Errorf("查询优惠券失败: %w", err)
		}

		// 2. 创建订单
		orderNo := fmt.Sprintf("%d%d", time.Now().UnixMilli(), task.UserID%10000)
		order := &orderModel.Order{
			OrderNo:     orderNo,
			UserID:      task.UserID,
			CouponID:    task.CouponID,
			TotalAmount: coupon.Value,
			PayAmount:   0,
			Status:      1,
			PaidAt:      func() *time.Time { t := time.Now(); return &t }(),
		}
		if err := tx.Create(order).Error; err != nil {
			return fmt.Errorf("create order: %w", err)
		}

		// 3. 记录用户领券
		userCoupon := &model.UserCoupon{
			UserID:   task.UserID,
			CouponID: task.CouponID,
			OrderID:  order.ID,
			Status:   0,
		}
		if err := tx.Create(userCoupon).Error; err != nil {
			return fmt.Errorf("create user_coupon: %w", err)
		}

		// 4. 更新已用数量
		if err := tx.Model(&model.Coupon{}).
			Where("id = ?", task.CouponID).
			UpdateColumn("used_count", gorm.Expr("used_count + 1")).Error; err != nil {
			return fmt.Errorf("update used_count: %w", err)
		}

		return nil
	})
}
