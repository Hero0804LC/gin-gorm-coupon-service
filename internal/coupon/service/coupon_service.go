package service

import (
	"context"
	"fmt"
	"time"

	"gin-gorm-coupon-service/internal/coupon/model"
	"gin-gorm-coupon-service/internal/coupon/repository"
)

type CouponService struct {
	couponRepo     repository.CouponRepo
	userCouponRepo repository.UserCouponRepo
}

func NewCouponService(couponRepo repository.CouponRepo, userCouponRepo repository.UserCouponRepo) *CouponService {
	return &CouponService{
		couponRepo:     couponRepo,
		userCouponRepo: userCouponRepo,
	}
}

// CreateCouponRequest 请求结构体
type CreateCouponRequest struct {
	Name       string  `json:"name" binding:"required"`
	Type       int8    `json:"type" binding:"required,oneof=1 2 3"`
	Amount     float64 `json:"amount"`
	MinAmount  float64 `json:"min_amount"`
	TotalCount int     `json:"total_count" binding:"required,min=1"`
	PerLimit   int     `json:"per_limit" binding:"required,min=1"`
	StartTime  string  `json:"start_time" binding:"required"`
	EndTime    string  `json:"end_time" binding:"required"`
}

// Create 创建优惠券（手动解析时间，避免 time.Time 零值写 DB 报 500）
func (s *CouponService) Create(ctx context.Context, req *CreateCouponRequest) error {

	// 解析时间（支持 "2026-01-01 00:00:00" 和 RFC3339 两种格式）
	startTime, err := parseTime(req.StartTime)
	if err != nil {
		return fmt.Errorf("start_time 格式错误: %w", err)
	}
	endTime, err := parseTime(req.EndTime)
	if err != nil {
		return fmt.Errorf("end_time 格式错误: %w", err)
	}

	coupon := &model.Coupon{
		Name:       req.Name,
		Type:       req.Type,
		Value:      req.Amount,
		MinAmount:  req.MinAmount,
		TotalCount: req.TotalCount,
		PerLimit:   req.PerLimit,
		StartTime:  startTime,
		EndTime:    endTime,
		Status:     1,
	}

	return s.couponRepo.Create(ctx, coupon)
}

// parseTime 支持多种时间格式解析
func parseTime(s string) (time.Time, error) {
	// 尝试 RFC3339
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	// 尝试 "2006-01-02 15:04:05"
	if t, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local); err == nil {
		return t, nil
	}
	// 尝试 "2006-01-02"
	if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("不支持的时间格式: %s（请用 2006-01-01 00:00:00 或 2006-01-01T00:00:00Z）", s)
}

// List 列表
func (s *CouponService) List(ctx context.Context) ([]*model.Coupon, error) {
	return s.couponRepo.List(ctx)
}

// GetByID 详情
func (s *CouponService) GetByID(ctx context.Context, id int64) (*model.Coupon, error) {
	return s.couponRepo.GetByID(ctx, id)
}
