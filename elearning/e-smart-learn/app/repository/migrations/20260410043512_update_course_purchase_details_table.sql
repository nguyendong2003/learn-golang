-- +goose Up
-- +goose StatementBegin
DROP VIEW IF EXISTS v_course_purchase_revenue_base;

ALTER TABLE course_purchase_details
    DROP COLUMN IF EXISTS price,
    ADD COLUMN price_original BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN price_final BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN coupon_id UUID REFERENCES coupons(id) ON DELETE SET NULL;

ALTER TABLE course_purchases
    DROP COLUMN IF EXISTS coupon_id,
    DROP COLUMN IF EXISTS amount,
    ADD COLUMN total_amount BIGINT NOT NULL DEFAULT 0;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE course_purchase_details
    DROP COLUMN IF EXISTS price_original,
    DROP COLUMN IF EXISTS price_final,
    DROP COLUMN IF EXISTS coupon_id,
    ADD COLUMN price BIGINT NOT NULL DEFAULT 0;

ALTER TABLE course_purchases
    ADD COLUMN coupon_id UUID REFERENCES coupons(id) ON DELETE SET NULL,
    ADD COLUMN amount BIGINT NOT NULL DEFAULT 0,
    DROP COLUMN IF EXISTS total_amount;
-- +goose StatementEnd
