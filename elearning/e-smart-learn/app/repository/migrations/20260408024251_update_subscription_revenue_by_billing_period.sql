-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fn_subscription_revenue_by_instructor(
	p_instructor_id UUID,
	p_from TIMESTAMPTZ DEFAULT NULL,
	p_to TIMESTAMPTZ DEFAULT NULL
)
RETURNS TABLE(
	payment_id UUID,
	subscription_id UUID,
	subscriber_user_id UUID,
	paid_at TIMESTAMPTZ,
	payment_amount BIGINT,
	payment_stripe_fee BIGINT,
	total_active_courses INTEGER,
	instructor_active_courses INTEGER,
	allocated_amount NUMERIC
)
LANGUAGE sql
AS $$
WITH successful_payments AS (
	SELECT
		p.id AS payment_id,
		p.subscription_id,
		p.amount,
		p.stripe_fee,
		p.paid_at AS paid_at,
		COALESCE(p.billing_period_start, p.paid_at) AS period_start,
		COALESCE(p.billing_period_end, p.paid_at) AS period_end,
		s.user_id
	FROM payments p
	JOIN subscriptions s ON s.id = p.subscription_id
	WHERE p.status = 'succeeded'
	  AND p.paid_at IS NOT NULL
	  AND (p_from IS NULL OR p.paid_at >= p_from)
	  AND (p_to IS NULL OR p.paid_at < p_to)
),
active_enrollments AS (
	SELECT
		sp.payment_id,
		sp.subscription_id,
		sp.user_id,
		sp.amount,
		sp.stripe_fee,
		sp.paid_at,
		e.course_id,
		c.user_id AS instructor_id,
		ROW_NUMBER() OVER (
			PARTITION BY sp.payment_id, e.course_id
			ORDER BY e.enrolled_at DESC, e.id DESC
		) AS rn
	FROM successful_payments sp
	JOIN enrollments e
		ON e.user_id = sp.user_id
		AND e.type = 'subscription'
		AND e.enrolled_at < sp.period_end
		AND (e.canceled_at IS NULL OR e.canceled_at > sp.period_start)
	JOIN courses c ON c.id = e.course_id
),
dedup_active_enrollments AS (
	SELECT
		payment_id,
		subscription_id,
		user_id,
		amount,
		stripe_fee,
		paid_at,
		course_id,
		instructor_id
	FROM active_enrollments
	WHERE rn = 1
),
total_course_count_per_payment AS (
	SELECT payment_id, COUNT(*)::INTEGER AS total_active_courses
	FROM dedup_active_enrollments
	GROUP BY payment_id
),
instructor_course_count_per_payment AS (
	SELECT
		payment_id,
		instructor_id,
		COUNT(*)::INTEGER AS instructor_active_courses
	FROM dedup_active_enrollments
	GROUP BY payment_id, instructor_id
)
SELECT
	ic.payment_id,
	d.subscription_id,
	d.user_id AS subscriber_user_id,
	d.paid_at,
	d.amount AS payment_amount,
	d.stripe_fee AS payment_stripe_fee,
	tc.total_active_courses,
	ic.instructor_active_courses,
	(d.amount::NUMERIC * ic.instructor_active_courses::NUMERIC / NULLIF(tc.total_active_courses::NUMERIC, 0)) AS allocated_amount
FROM instructor_course_count_per_payment ic
JOIN total_course_count_per_payment tc ON tc.payment_id = ic.payment_id
JOIN (
	SELECT DISTINCT payment_id, subscription_id, user_id, paid_at, amount, stripe_fee
	FROM dedup_active_enrollments
) d ON d.payment_id = ic.payment_id
WHERE ic.instructor_id = p_instructor_id
ORDER BY d.paid_at, ic.payment_id;
$$;

CREATE OR REPLACE FUNCTION get_admin_subscription_revenue_range(
	p_from TIMESTAMPTZ,
	p_to TIMESTAMPTZ
)
RETURNS TABLE (
	total_amount BIGINT,
	instructor_gross BIGINT,
	platform_gross BIGINT,
	stripe_fee BIGINT,
	instructor_net BIGINT,
	platform_net BIGINT
)
LANGUAGE sql
AS $$
WITH successful_payments AS (
	SELECT
		p.id AS payment_id,
		p.amount,
		p.stripe_fee,
		COALESCE(p.billing_period_start, p.paid_at) AS period_start,
		COALESCE(p.billing_period_end, p.paid_at) AS period_end,
		s.user_id
	FROM payments p
	JOIN subscriptions s ON s.id = p.subscription_id
	WHERE p.status = 'succeeded'
	  AND p.paid_at IS NOT NULL
	  AND p.paid_at >= p_from
	  AND p.paid_at < p_to
),
active_enrollments AS (
	SELECT
		sp.payment_id,
		e.course_id,
		ROW_NUMBER() OVER (
			PARTITION BY sp.payment_id, e.course_id
			ORDER BY e.enrolled_at DESC, e.id DESC
		) AS rn
	FROM successful_payments sp
	JOIN enrollments e
		ON e.user_id = sp.user_id
		AND e.type = 'subscription'
		AND e.enrolled_at < sp.period_end
		AND (e.canceled_at IS NULL OR e.canceled_at > sp.period_start)
),
eligible_payments AS (
	SELECT DISTINCT
		sp.payment_id,
		sp.amount,
		sp.stripe_fee
	FROM successful_payments sp
	JOIN active_enrollments ae ON ae.payment_id = sp.payment_id
	WHERE ae.rn = 1
)
SELECT
	total,
	igross,
	total - igross,
	fee,
	inet,
	(total - fee) - inet
FROM (
	SELECT
		total,
		fee,
		(total * 70 / 100) AS igross,
		((total - fee) * 70 / 100) AS inet
	FROM (
		SELECT
			COALESCE(SUM(ep.amount), 0)::BIGINT AS total,
			COALESCE(SUM(ep.stripe_fee), 0)::BIGINT AS fee
		FROM eligible_payments ep
	) base
) final;
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Down migration intentionally keeps current definitions.
-- If rollback is needed, re-run the previous subscription revenue migration.
SELECT 1;
-- +goose StatementEnd
