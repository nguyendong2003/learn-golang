-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_cp_paid_purchased_at
ON course_purchases(purchased_at)
WHERE status = 'paid';

-- =========================================
-- 1. Function: calculate price after discount
-- =========================================
CREATE OR REPLACE FUNCTION calculate_price_after_discount(
    price BIGINT,
    discount_type discount_type_enum,
    discount_value BIGINT
)
RETURNS NUMERIC
LANGUAGE plpgsql
AS $$
BEGIN
    IF discount_type IS NULL OR discount_value IS NULL THEN
        RETURN price;
    END IF;

    IF discount_type = 'percent' THEN
        RETURN GREATEST(price * (100 - discount_value) / 100.0, 0);

    ELSIF discount_type = 'amount' THEN
        RETURN GREATEST(price - discount_value, 0);

    ELSE
        RETURN price;
    END IF;
END;
$$;

-- =========================================
-- 2. View: base revenue data
-- =========================================
CREATE VIEW v_course_purchase_revenue_base AS
SELECT  
    cp.id AS course_purchase_id,
    cp.amount, 
    cp.stripe_fee,
    cp.currency,
    cp.purchased_at,

    cps.discount_type,
    cps.discount_value,

    cpd.course_id,
    cpd.price AS course_price,

    calculate_price_after_discount(
        cpd.price,
        cps.discount_type,
        cps.discount_value
    ) AS course_price_after_discount,

    c.user_id

FROM course_purchases cp 
LEFT JOIN coupons cps ON cps.id = cp.coupon_id 
JOIN course_purchase_details cpd ON cp.id = cpd.course_purchase_id
JOIN courses c ON c.id = cpd.course_id 
WHERE cp.status = 'paid';

-- =========================================
-- 3. CORE FUNCTION (TIMESTAMPTZ)
-- =========================================
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
				COALESCE(SUM(cp_amount), 0)::BIGINT AS total,

                COALESCE(SUM(cp_stripe_fee), 0)::BIGINT AS fee

            FROM (
                SELECT
                    course_purchase_id,
                    MAX(amount) AS cp_amount,
                    MAX(vcpb.stripe_fee) AS cp_stripe_fee
                FROM v_course_purchase_revenue_base vcpb
                WHERE user_id = p_user_id
                  AND purchased_at >= p_from
                  AND purchased_at < p_to
                GROUP BY course_purchase_id
            ) t
        ) base
    ) final;
END;
$$;

-- =========================================
-- 4. WRAPPER: TOTAL ALL TIME
-- =========================================
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

-- =========================================
-- 5. WRAPPER: BY DAY
-- =========================================
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

-- =========================================
-- 6. WRAPPER: BY MONTH
-- =========================================
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

-- =========================================
-- 7. WRAPPER: BY YEAR
-- =========================================
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

-- =========================================
-- 8. CORE FUNCTION: admin summary by time range (all instructors)
-- =========================================
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
                COALESCE(SUM(cp_amount), 0)::BIGINT AS total,
                COALESCE(SUM(cp_stripe_fee), 0)::BIGINT AS fee
            FROM (
                SELECT
                    cp.id AS course_purchase_id,
                    cp.amount AS cp_amount,
                    cp.stripe_fee AS cp_stripe_fee
                FROM course_purchases cp
                WHERE cp.status = 'paid'
                  AND cp.purchased_at >= p_from
                  AND cp.purchased_at < p_to
            ) t
        ) base
    ) final;
END;
$$;

-- =========================================
-- 9. WRAPPER: ADMIN TOTAL ALL TIME
-- =========================================
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

-- =========================================
-- 10. WRAPPER: ADMIN BY DAY
-- =========================================
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

-- =========================================
-- 11. WRAPPER: ADMIN BY MONTH
-- =========================================
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

-- =========================================
-- 12. WRAPPER: ADMIN BY YEAR
-- =========================================
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
DROP FUNCTION IF EXISTS get_instructor_course_purchase_revenue(UUID);
DROP FUNCTION IF EXISTS get_instructor_course_purchase_revenue_range(UUID, TIMESTAMPTZ, TIMESTAMPTZ);
DROP FUNCTION IF EXISTS get_instructor_course_purchase_revenue_by_day(UUID, DATE);
DROP FUNCTION IF EXISTS get_instructor_course_purchase_revenue_by_month(UUID, INT, INT);
DROP FUNCTION IF EXISTS get_instructor_course_purchase_revenue_by_year(UUID, INT);

DROP FUNCTION IF EXISTS get_admin_course_purchase_revenue();
DROP FUNCTION IF EXISTS get_admin_course_purchase_revenue_range(TIMESTAMPTZ, TIMESTAMPTZ);
DROP FUNCTION IF EXISTS get_admin_course_purchase_revenue_by_day(DATE);
DROP FUNCTION IF EXISTS get_admin_course_purchase_revenue_by_month(INT, INT);
DROP FUNCTION IF EXISTS get_admin_course_purchase_revenue_by_year(INT);

DROP VIEW IF EXISTS v_course_purchase_revenue_base;

DROP FUNCTION IF EXISTS calculate_price_after_discount(BIGINT, discount_type_enum, BIGINT);
-- +goose StatementEnd