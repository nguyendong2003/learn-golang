package repository

import (
	"context"
	"elearning-api/model"

	"github.com/google/uuid"
)

type CoursePurchaseRevenueShareRepository interface {
	Repository[model.CoursePurchaseRevenueShare]
	RebuildByCoursePurchaseID(ctx context.Context, purchaseID uuid.UUID) error
}

type coursePurchaseRevenueShareRepository struct {
	*repository[model.CoursePurchaseRevenueShare]
	db DbRepository
}

func NewCoursePurchaseRevenueShareRepository(db DbRepository) CoursePurchaseRevenueShareRepository {
	return &coursePurchaseRevenueShareRepository{
		repository: NewBaseRepository[model.CoursePurchaseRevenueShare](db),
		db:         db,
	}
}

func (r *coursePurchaseRevenueShareRepository) RebuildByCoursePurchaseID(ctx context.Context, purchaseID uuid.UUID) error {
	if err := r.db.GetDB().WithContext(ctx).
		Exec("DELETE FROM course_purchase_revenue_shares WHERE course_purchase_id = ?", purchaseID).Error; err != nil {
		return err
	}

	query := `
INSERT INTO course_purchase_revenue_shares (
	course_purchase_id,
	course_purchase_detail_id,
	purchaser_user_id,
	instructor_user_id,
	course_id,
	purchase_amount,
	purchase_stripe_fee,
	detail_amount,
	allocated_stripe_fee,
	instructor_gross,
	platform_gross,
	instructor_net,
	platform_net,
	purchased_at,
	created_at,
	updated_at
)
SELECT
	cp.id,
	cpd.id,
	cp.user_id,
	c.user_id,
	c.id,
	cp.amount,
	cp.stripe_fee,
	cpd.price,
	ROUND(cp.stripe_fee::numeric * cpd.price::numeric / NULLIF(cp.amount::numeric, 0))::bigint AS allocated_stripe_fee,
	ROUND(cpd.price::numeric * 0.7)::bigint AS instructor_gross,
	(cpd.price - ROUND(cpd.price::numeric * 0.7)::bigint) AS platform_gross,
	(ROUND(cpd.price::numeric * 0.7)::bigint - ROUND((cp.stripe_fee::numeric * cpd.price::numeric / NULLIF(cp.amount::numeric, 0)) * 0.7)::bigint) AS instructor_net,
	((cpd.price - ROUND(cpd.price::numeric * 0.7)::bigint) - (ROUND(cp.stripe_fee::numeric * cpd.price::numeric / NULLIF(cp.amount::numeric, 0))::bigint - ROUND((cp.stripe_fee::numeric * cpd.price::numeric / NULLIF(cp.amount::numeric, 0)) * 0.7)::bigint)) AS platform_net,
	COALESCE(cp.purchased_at, cp.created_at),
	NOW(),
	NOW()
FROM course_purchases cp
JOIN course_purchase_details cpd ON cpd.course_purchase_id = cp.id AND cpd.deleted_at IS NULL
JOIN courses c ON c.id = cpd.course_id AND c.deleted_at IS NULL
WHERE cp.id = ?
  AND cp.status = 'paid'
  AND cp.deleted_at IS NULL
`

	return r.db.GetDB().WithContext(ctx).Exec(query, purchaseID).Error
}
