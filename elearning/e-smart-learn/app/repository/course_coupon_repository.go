package repository

import (
	"context"
	"elearning-api/model"
	"time"

	"github.com/google/uuid"
)

type CourseCouponRepository interface {
	Repository[model.CourseCoupon]
	GetByCourseAndCouponID(ctx context.Context, courseID, couponID uuid.UUID, preloads []Preload) (*model.CourseCoupon, error)
	ListByCourseID(ctx context.Context, courseID uuid.UUID, preloads []Preload) ([]*model.CourseCoupon, error)
	GetDefaultByCourseID(ctx context.Context, courseID uuid.UUID, preloads []Preload) (*model.CourseCoupon, error)
	GetByCourseAndCouponCode(ctx context.Context, courseID uuid.UUID, code string, preloads []Preload) (*model.CourseCoupon, error)
	DeleteByCourseID(ctx context.Context, courseID uuid.UUID) error
	DeleteByCouponID(ctx context.Context, courseID, couponID uuid.UUID) error
	DeleteUnusableCoupons(ctx context.Context, now time.Time) (int64, error)
}

type courseCouponRepository struct {
	*repository[model.CourseCoupon]
}

func NewCourseCouponRepository(db DbRepository) CourseCouponRepository {
	return &courseCouponRepository{
		repository: NewBaseRepository[model.CourseCoupon](db),
	}
}

func (r *courseCouponRepository) GetByCourseAndCouponID(ctx context.Context, courseID, couponID uuid.UUID, preloads []Preload) (*model.CourseCoupon, error) {
	return r.Find(ctx, "course_id = ? AND coupon_id = ?", preloads, courseID, couponID)
}

func (r *courseCouponRepository) ListByCourseID(ctx context.Context, courseID uuid.UUID, preloads []Preload) ([]*model.CourseCoupon, error) {
	return r.FindAll(ctx, "course_id = ?", preloads, courseID)
}

func (r *courseCouponRepository) GetDefaultByCourseID(ctx context.Context, courseID uuid.UUID, preloads []Preload) (*model.CourseCoupon, error) {
	return r.Find(ctx, "course_id = ? AND is_default = ?", preloads, courseID, true)
}

func (r *courseCouponRepository) GetByCourseAndCouponCode(ctx context.Context, courseID uuid.UUID, code string, preloads []Preload) (*model.CourseCoupon, error) {
	db := r.baseQuery(ctx).
		Joins("JOIN coupons ON coupons.id = course_coupons.coupon_id AND coupons.deleted_at IS NULL").
		Where("course_coupons.course_id = ?", courseID).
		Where("coupons.code = ?", code)

	if len(preloads) > 0 {
		db = applyPreloads(db, preloads)
	}

	var courseCoupon model.CourseCoupon
	if err := db.First(&courseCoupon).Error; err != nil {
		return nil, err
	}

	return &courseCoupon, nil
}

func (r *courseCouponRepository) DeleteByCourseID(ctx context.Context, courseID uuid.UUID) error {
	return r.db.GetDB().WithContext(ctx).
		Where("course_id = ?", courseID).
		Delete(&model.CourseCoupon{}).Error
}

func (r *courseCouponRepository) DeleteByCouponID(ctx context.Context, courseID, couponID uuid.UUID) error {
	return r.db.GetDB().WithContext(ctx).
		Where("course_id = ? AND coupon_id = ?", courseID, couponID).
		Delete(&model.CourseCoupon{}).Error
}

func (r *courseCouponRepository) DeleteUnusableCoupons(ctx context.Context, now time.Time) (int64, error) {
	subQuery := r.db.GetDB().WithContext(ctx).
		Table("course_coupons AS cc").
		Joins("JOIN coupons c ON c.id = cc.coupon_id").
		Where("cc.deleted_at IS NULL").
		Where(
			"c.deleted_at IS NOT NULL OR c.is_active = ? OR (c.expires_at IS NOT NULL AND c.expires_at < ?) OR (c.max_redemptions IS NOT NULL AND c.max_redemptions > 0 AND c.current_redemptions >= c.max_redemptions)",
			false,
			now.UTC(),
		).
		Select("cc.id")

	result := r.db.GetDB().WithContext(ctx).
		Where("id IN (?)", subQuery).
		Delete(&model.CourseCoupon{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}
