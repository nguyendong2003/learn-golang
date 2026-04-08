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
	WITH purchase_ctx AS (
		SELECT
			cp.id AS course_purchase_id,
			cp.user_id AS purchaser_user_id,
			cp.amount AS purchase_amount,
			cp.stripe_fee AS purchase_stripe_fee,
			COALESCE(cp.purchased_at, cp.created_at) AS purchased_at,
			COALESCE(SUM(cpd.price), 0)::bigint AS gross_amount
		FROM course_purchases cp
		JOIN course_purchase_details cpd ON cpd.course_purchase_id = cp.id AND cpd.deleted_at IS NULL
		WHERE cp.id = ?
		  AND cp.status = 'paid'
		  AND cp.deleted_at IS NULL
		GROUP BY cp.id, cp.user_id, cp.amount, cp.stripe_fee, COALESCE(cp.purchased_at, cp.created_at)
	),
	detail_base AS (
		SELECT
			pc.course_purchase_id,
			cpd.id AS course_purchase_detail_id,
			pc.purchaser_user_id,
			c.user_id AS instructor_user_id,
			c.id AS course_id,
			pc.purchase_amount,
			pc.purchase_stripe_fee,
			cpd.price AS detail_amount,
			pc.purchased_at,
			pc.gross_amount,
			GREATEST(pc.gross_amount - pc.purchase_amount, 0)::bigint AS total_coupon_discount
		FROM purchase_ctx pc
		JOIN course_purchase_details cpd ON cpd.course_purchase_id = pc.course_purchase_id AND cpd.deleted_at IS NULL
		JOIN courses c ON c.id = cpd.course_id AND c.deleted_at IS NULL
	),
	detail_calc AS (
		SELECT
			db.*,
			CASE
				WHEN db.total_coupon_discount <= 0 OR db.gross_amount <= 0 THEN 0::bigint
				ELSE ROUND(db.total_coupon_discount::numeric * db.detail_amount::numeric / db.gross_amount::numeric)::bigint
			END AS detail_coupon_discount
		FROM detail_base db
	),
	detail_alloc AS (
		SELECT
			dc.course_purchase_id,
			dc.course_purchase_detail_id,
			dc.purchaser_user_id,
			dc.instructor_user_id,
			dc.course_id,
			dc.purchase_amount,
			dc.purchase_stripe_fee,
			dc.detail_amount,
			dc.detail_coupon_discount,
			GREATEST(dc.detail_amount - dc.detail_coupon_discount, 0)::bigint AS detail_net_amount,
			CASE
				WHEN dc.purchase_amount <= 0 THEN 0::bigint
				ELSE ROUND(dc.purchase_stripe_fee::numeric * GREATEST(dc.detail_amount - dc.detail_coupon_discount, 0)::numeric / dc.purchase_amount::numeric)::bigint
			END AS allocated_stripe_fee,
			dc.purchased_at
		FROM detail_calc dc
	)
	INSERT INTO course_purchase_revenue_shares (
		course_purchase_id,
		course_purchase_detail_id,
		purchaser_user_id,
		instructor_user_id,
		course_id,
		purchase_amount,
		purchase_stripe_fee,
		detail_amount,
		detail_coupon_discount,
		detail_net_amount,
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
		da.course_purchase_id,
		da.course_purchase_detail_id,
		da.purchaser_user_id,
		da.instructor_user_id,
		da.course_id,
		da.purchase_amount,
		da.purchase_stripe_fee,
		da.detail_amount,
		da.detail_coupon_discount,
		da.detail_net_amount,
		da.allocated_stripe_fee,
		ROUND(da.detail_net_amount::numeric * 0.7)::bigint AS instructor_gross,
		(da.detail_net_amount - ROUND(da.detail_net_amount::numeric * 0.7)::bigint) AS platform_gross,
		(ROUND(da.detail_net_amount::numeric * 0.7)::bigint - ROUND(da.allocated_stripe_fee::numeric * 0.7)::bigint) AS instructor_net,
		((da.detail_net_amount - ROUND(da.detail_net_amount::numeric * 0.7)::bigint) - (da.allocated_stripe_fee - ROUND(da.allocated_stripe_fee::numeric * 0.7)::bigint)) AS platform_net,
		da.purchased_at,
		NOW(),
		NOW()
	FROM detail_alloc da
	`

	return r.db.GetDB().WithContext(ctx).Exec(query, purchaseID).Error
}
