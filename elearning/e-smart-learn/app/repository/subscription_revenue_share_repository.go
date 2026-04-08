package repository

import (
	"context"
	"elearning-api/model"

	"github.com/google/uuid"
)

type SubscriptionRevenueShareRepository interface {
	Repository[model.SubscriptionRevenueShare]
	RebuildByPaymentID(ctx context.Context, paymentID uuid.UUID) error
}

type subscriptionRevenueShareRepository struct {
	*repository[model.SubscriptionRevenueShare]
	db DbRepository
}

func NewSubscriptionRevenueShareRepository(db DbRepository) SubscriptionRevenueShareRepository {
	return &subscriptionRevenueShareRepository{
		repository: NewBaseRepository[model.SubscriptionRevenueShare](db),
		db:         db,
	}
}

func (r *subscriptionRevenueShareRepository) RebuildByPaymentID(ctx context.Context, paymentID uuid.UUID) error {
	// Keep ledger idempotent: rebuild all rows for this payment from scratch.
	if err := r.db.GetDB().WithContext(ctx).
		Exec("DELETE FROM subscription_revenue_shares WHERE payment_id = ?", paymentID).Error; err != nil {
		return err
	}

	query := `
WITH payment_ctx AS (
	SELECT
		p.id AS payment_id,
		p.subscription_id,
		s.user_id AS subscriber_user_id,
		p.amount,
		p.stripe_fee,
		(COALESCE(p.billing_period_start, p.paid_at) AT TIME ZONE 'UTC') AS period_start,
		(COALESCE(p.billing_period_end, p.paid_at) AT TIME ZONE 'UTC') AS period_end,
		(p.paid_at AT TIME ZONE 'UTC') AS paid_at
	FROM payments p
	JOIN subscriptions s ON s.id = p.subscription_id
	WHERE p.id = ?
	  AND p.status = 'succeeded'
	  AND p.paid_at IS NOT NULL
),
active_enrollments AS (
	SELECT
		pc.payment_id,
		pc.subscription_id,
		pc.subscriber_user_id,
		pc.amount,
		pc.stripe_fee,
		pc.period_start,
		pc.period_end,
		pc.paid_at,
		e.course_id,
		c.user_id AS instructor_user_id,
		ROW_NUMBER() OVER (
			PARTITION BY pc.payment_id, e.course_id
			ORDER BY e.enrolled_at DESC, e.id DESC
		) AS rn
	FROM payment_ctx pc
	JOIN enrollments e
		ON e.user_id = pc.subscriber_user_id
		AND e.type = 'subscription'
		AND e.enrolled_at < pc.period_end
		AND (e.canceled_at IS NULL OR e.canceled_at > pc.period_start)
	JOIN courses c ON c.id = e.course_id
),
dedup AS (
	SELECT *
	FROM active_enrollments
	WHERE rn = 1
),
totals AS (
	SELECT payment_id, COUNT(*)::int AS total_active_courses
	FROM dedup
	GROUP BY payment_id
),
instructor_counts AS (
	SELECT
		payment_id,
		subscription_id,
		subscriber_user_id,
		instructor_user_id,
		amount,
		stripe_fee,
		period_start,
		period_end,
		paid_at,
		COUNT(*)::int AS instructor_active_courses
	FROM dedup
	GROUP BY
		payment_id,
		subscription_id,
		subscriber_user_id,
		instructor_user_id,
		amount,
		stripe_fee,
		period_start,
		period_end,
		paid_at
)
INSERT INTO subscription_revenue_shares (
	payment_id,
	subscription_id,
	subscriber_user_id,
	instructor_user_id,
	total_active_courses,
	instructor_active_courses,
	payment_amount,
	payment_stripe_fee,
	allocated_amount,
	instructor_gross,
	platform_gross,
	allocated_stripe_fee,
	instructor_net,
	platform_net,
	billing_period_start,
	billing_period_end,
	paid_at,
	created_at,
	updated_at
)
SELECT
	ic.payment_id,
	ic.subscription_id,
	ic.subscriber_user_id,
	ic.instructor_user_id,
	t.total_active_courses,
	ic.instructor_active_courses,
	ic.amount,
	ic.stripe_fee,
	ROUND(ic.amount::numeric * ic.instructor_active_courses::numeric / NULLIF(t.total_active_courses::numeric, 0))::bigint AS allocated_amount,
	ROUND((ic.amount::numeric * ic.instructor_active_courses::numeric / NULLIF(t.total_active_courses::numeric, 0)) * 0.7)::bigint AS instructor_gross,
	(ROUND(ic.amount::numeric * ic.instructor_active_courses::numeric / NULLIF(t.total_active_courses::numeric, 0))::bigint
	 - ROUND((ic.amount::numeric * ic.instructor_active_courses::numeric / NULLIF(t.total_active_courses::numeric, 0)) * 0.7)::bigint) AS platform_gross,
	ROUND(ic.stripe_fee::numeric * ic.instructor_active_courses::numeric / NULLIF(t.total_active_courses::numeric, 0))::bigint AS allocated_stripe_fee,
	(ROUND((ic.amount::numeric * ic.instructor_active_courses::numeric / NULLIF(t.total_active_courses::numeric, 0)) * 0.7)::bigint
	 - ROUND((ic.stripe_fee::numeric * ic.instructor_active_courses::numeric / NULLIF(t.total_active_courses::numeric, 0)) * 0.7)::bigint) AS instructor_net,
	((ROUND(ic.amount::numeric * ic.instructor_active_courses::numeric / NULLIF(t.total_active_courses::numeric, 0))::bigint
	  - ROUND((ic.amount::numeric * ic.instructor_active_courses::numeric / NULLIF(t.total_active_courses::numeric, 0)) * 0.7)::bigint)
	 - (ROUND(ic.stripe_fee::numeric * ic.instructor_active_courses::numeric / NULLIF(t.total_active_courses::numeric, 0))::bigint
	    - ROUND((ic.stripe_fee::numeric * ic.instructor_active_courses::numeric / NULLIF(t.total_active_courses::numeric, 0)) * 0.7)::bigint)) AS platform_net,
	ic.period_start,
	ic.period_end,
	ic.paid_at,
	NOW(),
	NOW()
FROM instructor_counts ic
JOIN totals t ON t.payment_id = ic.payment_id;
`

	return r.db.GetDB().WithContext(ctx).Exec(query, paymentID).Error
}
