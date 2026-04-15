-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE VIEW v_course_purchase_revenue_base AS
SELECT
    cp.id AS course_purchase_id,
    cp.total_amount,
    cp.stripe_fee,
    cp.currency,
    cp.purchased_at,
    cpd.coupon_id,
    cps.discount_type,
    cps.discount_value,
    cpd.course_id,
    cpd.price_original AS course_price,
    cpd.price_final AS course_price_after_discount,
    c.user_id
FROM course_purchases cp
JOIN course_purchase_details cpd ON cp.id = cpd.course_purchase_id
JOIN courses c ON c.id = cpd.course_id
LEFT JOIN coupons cps ON cps.id = cpd.coupon_id
WHERE cp.status = 'paid';

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
                COALESCE(SUM(cp_total_amount), 0)::BIGINT AS total,
                COALESCE(SUM(cp_stripe_fee), 0)::BIGINT AS fee
            FROM (
                SELECT
                    DISTINCT cp.id AS course_purchase_id,
                    cp.total_amount AS cp_total_amount,
                    cp.stripe_fee AS cp_stripe_fee
                FROM course_purchases cp
                JOIN course_purchase_details cpd ON cpd.course_purchase_id = cp.id AND cpd.deleted_at IS NULL
                JOIN courses c ON c.id = cpd.course_id AND c.deleted_at IS NULL
                WHERE c.user_id = p_user_id
                  AND cp.status = 'paid'
                  AND cp.deleted_at IS NULL
                  AND cp.purchased_at >= p_from
                  AND cp.purchased_at < p_to
            ) t
        ) base
    ) final;
END;
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
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT *
    FROM get_instructor_course_purchase_revenue_range(
        p_user_id,
        TIMESTAMPTZ '1970-01-01 00:00:00+00',
        NOW()
    );
END;
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
    FROM get_instructor_course_purchase_revenue_range(
        p_user_id,
        v_start,
        v_start + INTERVAL '1 day'
    );
END;
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
    FROM get_instructor_course_purchase_revenue_range(
        p_user_id,
        v_start,
        v_start + INTERVAL '1 month'
    );
END;
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
    FROM get_instructor_course_purchase_revenue_range(
        p_user_id,
        v_start,
        v_start + INTERVAL '1 year'
    );
END;
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
                COALESCE(SUM(cp.total_amount), 0)::BIGINT AS total,
                COALESCE(SUM(cp.stripe_fee), 0)::BIGINT AS fee
            FROM course_purchases cp
            WHERE cp.status = 'paid'
              AND cp.deleted_at IS NULL
              AND cp.purchased_at >= p_from
              AND cp.purchased_at < p_to
        ) base
    ) final;
END;
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
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT *
    FROM get_admin_course_purchase_revenue_range(
        TIMESTAMPTZ '1970-01-01 00:00:00+00',
        NOW()
    );
END;
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
    FROM get_admin_course_purchase_revenue_range(
        v_start,
        v_start + INTERVAL '1 day'
    );
END;
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
    FROM get_admin_course_purchase_revenue_range(
        v_start,
        v_start + INTERVAL '1 month'
    );
END;
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
    FROM get_admin_course_purchase_revenue_range(
        v_start,
        v_start + INTERVAL '1 year'
    );
END;
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- no-op: keep latest revenue function definitions
SELECT 1;
-- +goose StatementEnd
