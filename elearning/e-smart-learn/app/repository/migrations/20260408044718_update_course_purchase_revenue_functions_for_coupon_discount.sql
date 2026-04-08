-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION get_instructor_course_purchase_revenue_range(
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
LANGUAGE sql
AS $$
SELECT
	COALESCE(SUM(detail_net_amount), 0)::BIGINT AS total_amount,
	COALESCE(SUM(instructor_gross), 0)::BIGINT AS instructor_gross,
	COALESCE(SUM(platform_gross), 0)::BIGINT AS platform_gross,
	COALESCE(SUM(allocated_stripe_fee), 0)::BIGINT AS stripe_fee,
	COALESCE(SUM(instructor_net), 0)::BIGINT AS instructor_net,
	COALESCE(SUM(platform_net), 0)::BIGINT AS platform_net
FROM course_purchase_revenue_shares
WHERE instructor_user_id = p_user_id
  AND (purchased_at IS NOT NULL)
  AND purchased_at >= p_from
  AND purchased_at < p_to
  AND deleted_at IS NULL;
$$;

CREATE OR REPLACE FUNCTION get_instructor_course_purchase_revenue(p_user_id UUID)
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
SELECT * FROM get_instructor_course_purchase_revenue_range(
	p_user_id,
	TIMESTAMPTZ '1970-01-01 00:00:00+00',
	NOW()
);
$$;

CREATE OR REPLACE FUNCTION get_instructor_course_purchase_revenue_by_day(
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
LANGUAGE sql
AS $$
SELECT
	p_date,
	r.*
FROM get_instructor_course_purchase_revenue_range(
	p_user_id,
	p_date::TIMESTAMPTZ,
	(p_date::TIMESTAMPTZ + INTERVAL '1 day')
) r;
$$;

CREATE OR REPLACE FUNCTION get_instructor_course_purchase_revenue_by_month(
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
LANGUAGE sql
AS $$
SELECT
	TO_CHAR(MAKE_DATE(p_year, p_month, 1), 'YYYY-MM'),
	r.*
FROM get_instructor_course_purchase_revenue_range(
	p_user_id,
	MAKE_DATE(p_year, p_month, 1)::TIMESTAMPTZ,
	(MAKE_DATE(p_year, p_month, 1)::TIMESTAMPTZ + INTERVAL '1 month')
) r;
$$;

CREATE OR REPLACE FUNCTION get_instructor_course_purchase_revenue_by_year(
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
LANGUAGE sql
AS $$
SELECT
	p_year,
	r.*
FROM get_instructor_course_purchase_revenue_range(
	p_user_id,
	MAKE_DATE(p_year, 1, 1)::TIMESTAMPTZ,
	(MAKE_DATE(p_year, 1, 1)::TIMESTAMPTZ + INTERVAL '1 year')
) r;
$$;

CREATE OR REPLACE FUNCTION get_admin_course_purchase_revenue_range(
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
SELECT
	COALESCE(SUM(detail_net_amount), 0)::BIGINT AS total_amount,
	COALESCE(SUM(instructor_gross), 0)::BIGINT AS instructor_gross,
	COALESCE(SUM(platform_gross), 0)::BIGINT AS platform_gross,
	COALESCE(SUM(allocated_stripe_fee), 0)::BIGINT AS stripe_fee,
	COALESCE(SUM(instructor_net), 0)::BIGINT AS instructor_net,
	COALESCE(SUM(platform_net), 0)::BIGINT AS platform_net
FROM course_purchase_revenue_shares
WHERE purchased_at IS NOT NULL
  AND purchased_at >= p_from
  AND purchased_at < p_to
  AND deleted_at IS NULL;
$$;

CREATE OR REPLACE FUNCTION get_admin_course_purchase_revenue()
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
SELECT * FROM get_admin_course_purchase_revenue_range(
	TIMESTAMPTZ '1970-01-01 00:00:00+00',
	NOW()
);
$$;

CREATE OR REPLACE FUNCTION get_admin_course_purchase_revenue_by_day(
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
FROM get_admin_course_purchase_revenue_range(
	p_date::TIMESTAMPTZ,
	(p_date::TIMESTAMPTZ + INTERVAL '1 day')
) r;
$$;

CREATE OR REPLACE FUNCTION get_admin_course_purchase_revenue_by_month(
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
FROM get_admin_course_purchase_revenue_range(
	MAKE_DATE(p_year, p_month, 1)::TIMESTAMPTZ,
	(MAKE_DATE(p_year, p_month, 1)::TIMESTAMPTZ + INTERVAL '1 month')
) r;
$$;

CREATE OR REPLACE FUNCTION get_admin_course_purchase_revenue_by_year(
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
FROM get_admin_course_purchase_revenue_range(
	MAKE_DATE(p_year, 1, 1)::TIMESTAMPTZ,
	(MAKE_DATE(p_year, 1, 1)::TIMESTAMPTZ + INTERVAL '1 year')
) r;
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 1;
-- +goose StatementEnd
