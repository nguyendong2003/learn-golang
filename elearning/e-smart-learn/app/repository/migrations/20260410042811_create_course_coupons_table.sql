-- +goose Up
-- +goose StatementBegin
CREATE TABLE course_coupons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    coupon_id UUID NOT NULL REFERENCES coupons(id) ON DELETE CASCADE,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_course_coupons_course_id_coupon_id ON course_coupons(course_id, coupon_id);
CREATE INDEX idx_course_coupons_course_id ON course_coupons(course_id);
CREATE INDEX idx_course_coupons_coupon_id ON course_coupons(coupon_id);
CREATE INDEX idx_course_coupons_deleted_at ON course_coupons(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_course_coupons_deleted_at;
DROP INDEX IF EXISTS idx_course_coupons_coupon_id;
DROP INDEX IF EXISTS idx_course_coupons_course_id;
DROP INDEX IF EXISTS idx_course_coupons_course_id_coupon_id;
DROP TABLE IF EXISTS course_coupons;
-- +goose StatementEnd
