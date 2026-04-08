-- +goose Up
-- +goose StatementBegin
CREATE TYPE course_purchase_status_enum AS ENUM ('pending', 'paid', 'failed', 'refunded');

CREATE TABLE IF NOT EXISTS course_purchases (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  coupon_id UUID REFERENCES coupons(id) ON DELETE SET NULL,
  stripe_checkout_session_id VARCHAR(255) NOT NULL UNIQUE,
  stripe_payment_intent_id VARCHAR(255),
  amount BIGINT NOT NULL DEFAULT 0,
  currency VARCHAR(10) NOT NULL DEFAULT 'usd',
  status course_purchase_status_enum NOT NULL DEFAULT 'pending',
  purchased_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_course_purchases_user_id ON course_purchases(user_id);
CREATE INDEX IF NOT EXISTS idx_course_purchases_coupon_id ON course_purchases(coupon_id);
CREATE INDEX IF NOT EXISTS idx_course_purchases_payment_intent ON course_purchases(stripe_payment_intent_id);
CREATE INDEX IF NOT EXISTS idx_course_purchases_status ON course_purchases(status);
CREATE INDEX IF NOT EXISTS idx_course_purchases_deleted_at ON course_purchases(deleted_at);

CREATE TABLE IF NOT EXISTS course_purchase_details (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  course_purchase_id UUID NOT NULL REFERENCES course_purchases(id) ON DELETE CASCADE,
  course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
  price BIGINT NOT NULL DEFAULT 0,
  currency VARCHAR(10) NOT NULL DEFAULT 'usd',
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_course_purchase_details_course_purchase_id ON course_purchase_details(course_purchase_id);
CREATE INDEX IF NOT EXISTS idx_course_purchase_details_course_id ON course_purchase_details(course_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_course_purchase_details_purchase_course_unique ON course_purchase_details(course_purchase_id, course_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_course_purchase_details_purchase_course_unique;
DROP INDEX IF EXISTS idx_course_purchase_details_course_id;
DROP INDEX IF EXISTS idx_course_purchase_details_course_purchase_id;
DROP TABLE IF EXISTS course_purchase_details;
DROP INDEX IF EXISTS idx_course_purchases_deleted_at;
DROP INDEX IF EXISTS idx_course_purchases_status;
DROP INDEX IF EXISTS idx_course_purchases_payment_intent;
DROP INDEX IF EXISTS idx_course_purchases_coupon_id;
DROP INDEX IF EXISTS idx_course_purchases_user_id;
DROP TABLE IF EXISTS course_purchases;
DROP TYPE IF EXISTS course_purchase_status_enum;
-- +goose StatementEnd
