package repository

import (
	"context"
	"elearning-api/consts"
	"elearning-api/model"

	"github.com/google/uuid"
)

type CouponRepository interface {
	Repository[model.Coupon]
	GetByCode(ctx context.Context, code string, preloads []Preload) (*model.Coupon, error)
	CountPaidRedemptions(ctx context.Context, couponID uuid.UUID) (int64, error)
}

type couponRepository struct {
	*repository[model.Coupon]
}

func NewCouponRepository(db DbRepository) CouponRepository {
	return &couponRepository{
		repository: NewBaseRepository[model.Coupon](db),
	}
}

func (r *couponRepository) GetByCode(ctx context.Context, code string, preloads []Preload) (*model.Coupon, error) {
	return r.Find(ctx, "code = ?", preloads, code)
}

func (r *couponRepository) CountPaidRedemptions(ctx context.Context, couponID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.GetDB().WithContext(ctx).
		Table("course_purchase_details cpd").
		Joins("JOIN course_purchases cp ON cp.id = cpd.course_purchase_id").
		Where("cpd.coupon_id = ? AND cp.status = ? AND cp.deleted_at IS NULL AND cpd.deleted_at IS NULL", couponID, consts.CoursePurchaseStatusPaid).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}
