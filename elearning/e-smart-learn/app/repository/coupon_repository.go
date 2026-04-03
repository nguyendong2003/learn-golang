package repository

import (
	"context"
	"elearning-api/model"
)

type CouponRepository interface {
	Repository[model.Coupon]
	GetByCode(ctx context.Context, code string, preloads []Preload) (*model.Coupon, error)
}

type courseCouponRepository struct {
	*repository[model.Coupon]
}

func NewCouponRepository(db DbRepository) CouponRepository {
	return &courseCouponRepository{
		repository: NewBaseRepository[model.Coupon](db),
	}
}

func (r *courseCouponRepository) GetByCode(ctx context.Context, code string, preloads []Preload) (*model.Coupon, error) {
	return r.Find(ctx, "code = ?", preloads, code)
}
