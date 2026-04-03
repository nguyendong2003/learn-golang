-- +goose Up
-- +goose StatementBegin
-- =========================================
-- 1. DETAIL FUNCTION: allocate subscription payment to instructor
-- =========================================
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
		p.paid_at::TIMESTAMPTZ AS paid_at,
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
			-- ORDER BY COALESCE(e.enrolled_at, e.created_at) DESC, e.id DESC
		) AS rn
	FROM successful_payments sp
	JOIN enrollments e
	ON e.user_id = sp.user_id
	AND e.type = 'subscription'
	--  AND COALESCE(e.enrolled_at, e.created_at) <= sp.paid_at
    AND e.enrolled_at <= sp.paid_at
    AND (e.canceled_at IS NULL OR e.canceled_at > sp.paid_at)
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
	SELECT
		payment_id,
		COUNT(*)::INTEGER AS total_active_courses
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
	(
		d.amount::NUMERIC
		* ic.instructor_active_courses::NUMERIC
		/ NULLIF(tc.total_active_courses::NUMERIC, 0)
	) AS allocated_amount
FROM instructor_course_count_per_payment ic
JOIN total_course_count_per_payment tc ON tc.payment_id = ic.payment_id
JOIN (
	SELECT DISTINCT payment_id, subscription_id, user_id, paid_at, amount, stripe_fee
	FROM dedup_active_enrollments
) d ON d.payment_id = ic.payment_id
WHERE ic.instructor_id = p_instructor_id
ORDER BY d.paid_at, ic.payment_id;
$$;

-- =========================================
-- 2. CORE FUNCTION: summary by time range
-- =========================================
CREATE OR REPLACE FUNCTION get_instructor_subscription_revenue_range(
	p_user_id UUID,
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
LANGUAGE plpgsql
AS $$
BEGIN
	RETURN QUERY
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
				COALESCE(ROUND(SUM(allocated_amount)), 0)::BIGINT AS total,
				COALESCE(
					ROUND(SUM(
						payment_stripe_fee
						* (allocated_amount / NULLIF(payment_amount::NUMERIC, 0))
					)),
					0
				)::BIGINT AS fee
			FROM fn_subscription_revenue_by_instructor(p_user_id, p_from, p_to)
		) base
	) final;
END;
$$;

-- =========================================
-- 3. WRAPPER: TOTAL ALL TIME
-- =========================================
CREATE OR REPLACE FUNCTION get_instructor_subscription_revenue(p_user_id UUID)
RETURNS TABLE (
	total_amount BIGINT,
	instructor_gross BIGINT,
	platform_gross BIGINT,
	stripe_fee BIGINT,
	instructor_net BIGINT,
	platform_net BIGINT
)
LANGUAGE plpgsql
AS $$
BEGIN
	RETURN QUERY
	SELECT *
	FROM get_instructor_subscription_revenue_range(
		p_user_id,
		TIMESTAMPTZ '1970-01-01 00:00:00+00',
		NOW()
	);
END;
$$;

-- =========================================
-- 4. WRAPPER: BY DAY
-- =========================================
CREATE OR REPLACE FUNCTION get_instructor_subscription_revenue_by_day(
	p_user_id UUID,
	p_date DATE
)
RETURNS TABLE (
	revenue_date DATE,
	total_amount BIGINT,
	instructor_gross BIGINT,
	platform_gross BIGINT,
	stripe_fee BIGINT,
	instructor_net BIGINT,
	platform_net BIGINT
)
LANGUAGE plpgsql
AS $$
DECLARE
	v_start TIMESTAMPTZ;
BEGIN
	v_start := p_date::TIMESTAMPTZ;

	RETURN QUERY
	SELECT
		p_date,
		*
	FROM get_instructor_subscription_revenue_range(
		p_user_id,
		v_start,
		v_start + INTERVAL '1 day'
	);
END;
$$;

-- =========================================
-- 5. WRAPPER: BY MONTH
-- =========================================
CREATE OR REPLACE FUNCTION get_instructor_subscription_revenue_by_month(
	p_user_id UUID,
	p_year INT,
	p_month INT
)
RETURNS TABLE (
	revenue_month TEXT,
	total_amount BIGINT,
	instructor_gross BIGINT,
	platform_gross BIGINT,
	stripe_fee BIGINT,
	instructor_net BIGINT,
	platform_net BIGINT
)
LANGUAGE plpgsql
AS $$
DECLARE
	v_start TIMESTAMPTZ;
BEGIN
	v_start := MAKE_DATE(p_year, p_month, 1)::TIMESTAMPTZ;

	RETURN QUERY
	SELECT
		TO_CHAR(v_start::DATE, 'YYYY-MM'),
		*
	FROM get_instructor_subscription_revenue_range(
		p_user_id,
		v_start,
		v_start + INTERVAL '1 month'
	);
END;
$$;

-- =========================================
-- 6. WRAPPER: BY YEAR
-- =========================================
CREATE OR REPLACE FUNCTION get_instructor_subscription_revenue_by_year(
	p_user_id UUID,
	p_year INT
)
RETURNS TABLE (
	revenue_year INT,
	total_amount BIGINT,
	instructor_gross BIGINT,
	platform_gross BIGINT,
	stripe_fee BIGINT,
	instructor_net BIGINT,
	platform_net BIGINT
)
LANGUAGE plpgsql
AS $$
DECLARE
	v_start TIMESTAMPTZ;
BEGIN
	v_start := MAKE_DATE(p_year, 1, 1)::TIMESTAMPTZ;

	RETURN QUERY
	SELECT
		p_year,
		*
	FROM get_instructor_subscription_revenue_range(
		p_user_id,
		v_start,
		v_start + INTERVAL '1 year'
	);
END;
$$;

-- =========================================
-- 7. CORE FUNCTION: admin summary by time range (all instructors)
-- =========================================
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
		p.subscription_id,
		p.amount,
		p.stripe_fee,
		p.paid_at::TIMESTAMPTZ AS paid_at,
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
	AND e.enrolled_at <= sp.paid_at
	AND (e.canceled_at IS NULL OR e.canceled_at > sp.paid_at)
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

-- =========================================
-- 8. WRAPPER: ADMIN TOTAL ALL TIME
-- =========================================
CREATE OR REPLACE FUNCTION get_admin_subscription_revenue()
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
SELECT *
FROM get_admin_subscription_revenue_range(
	TIMESTAMPTZ '1970-01-01 00:00:00+00',
	NOW()
);
$$;

-- =========================================
-- 9. WRAPPER: ADMIN BY DAY
-- =========================================
CREATE OR REPLACE FUNCTION get_admin_subscription_revenue_by_day(
	p_date DATE
)
RETURNS TABLE (
	revenue_date DATE,
	total_amount BIGINT,
	instructor_gross BIGINT,
	platform_gross BIGINT,
	stripe_fee BIGINT,
	instructor_net BIGINT,
	platform_net BIGINT
)
LANGUAGE sql
AS $$
SELECT
	p_date,
	r.*
FROM get_admin_subscription_revenue_range(
	p_date::TIMESTAMPTZ,
	(p_date::TIMESTAMPTZ + INTERVAL '1 day')
) r;
$$;

-- =========================================
-- 10. WRAPPER: ADMIN BY MONTH
-- =========================================
CREATE OR REPLACE FUNCTION get_admin_subscription_revenue_by_month(
	p_year INT,
	p_month INT
)
RETURNS TABLE (
	revenue_month TEXT,
	total_amount BIGINT,
	instructor_gross BIGINT,
	platform_gross BIGINT,
	stripe_fee BIGINT,
	instructor_net BIGINT,
	platform_net BIGINT
)
LANGUAGE sql
AS $$
SELECT
	TO_CHAR(MAKE_DATE(p_year, p_month, 1), 'YYYY-MM'),
	r.*
FROM get_admin_subscription_revenue_range(
	MAKE_DATE(p_year, p_month, 1)::TIMESTAMPTZ,
	(MAKE_DATE(p_year, p_month, 1)::TIMESTAMPTZ + INTERVAL '1 month')
) r;
$$;

-- =========================================
-- 11. WRAPPER: ADMIN BY YEAR
-- =========================================
CREATE OR REPLACE FUNCTION get_admin_subscription_revenue_by_year(
	p_year INT
)
RETURNS TABLE (
	revenue_year INT,
	total_amount BIGINT,
	instructor_gross BIGINT,
	platform_gross BIGINT,
	stripe_fee BIGINT,
	instructor_net BIGINT,
	platform_net BIGINT
)
LANGUAGE sql
AS $$
SELECT
	p_year,
	r.*
FROM get_admin_subscription_revenue_range(
	MAKE_DATE(p_year, 1, 1)::TIMESTAMPTZ,
	(MAKE_DATE(p_year, 1, 1)::TIMESTAMPTZ + INTERVAL '1 year')
) r;
$$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS get_instructor_subscription_revenue(UUID);
DROP FUNCTION IF EXISTS get_instructor_subscription_revenue_range(UUID, TIMESTAMPTZ, TIMESTAMPTZ);
DROP FUNCTION IF EXISTS get_instructor_subscription_revenue_by_day(UUID, DATE);
DROP FUNCTION IF EXISTS get_instructor_subscription_revenue_by_month(UUID, INT, INT);
DROP FUNCTION IF EXISTS get_instructor_subscription_revenue_by_year(UUID, INT);
DROP FUNCTION IF EXISTS fn_subscription_revenue_by_instructor(UUID, TIMESTAMPTZ, TIMESTAMPTZ);

DROP FUNCTION IF EXISTS get_admin_subscription_revenue();
DROP FUNCTION IF EXISTS get_admin_subscription_revenue_range(TIMESTAMPTZ, TIMESTAMPTZ);
DROP FUNCTION IF EXISTS get_admin_subscription_revenue_by_day(DATE);
DROP FUNCTION IF EXISTS get_admin_subscription_revenue_by_month(INT, INT);
DROP FUNCTION IF EXISTS get_admin_subscription_revenue_by_year(INT);
-- +goose StatementEnd
