-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS course_purchase_revenue_shares (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  course_purchase_id UUID NOT NULL REFERENCES course_purchases(id) ON DELETE CASCADE,
  course_purchase_detail_id UUID NOT NULL REFERENCES course_purchase_details(id) ON DELETE CASCADE,
  purchaser_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  instructor_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
  purchase_amount BIGINT NOT NULL DEFAULT 0,
  purchase_stripe_fee BIGINT NOT NULL DEFAULT 0,
  detail_amount BIGINT NOT NULL DEFAULT 0,
  detail_coupon_discount BIGINT NOT NULL DEFAULT 0,
  detail_net_amount BIGINT NOT NULL DEFAULT 0,
  allocated_stripe_fee BIGINT NOT NULL DEFAULT 0,
  instructor_gross BIGINT NOT NULL DEFAULT 0,
  platform_gross BIGINT NOT NULL DEFAULT 0,
  instructor_net BIGINT NOT NULL DEFAULT 0,
  platform_net BIGINT NOT NULL DEFAULT 0,
  purchased_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_course_purchase_revenue_shares_detail_id
  ON course_purchase_revenue_shares(course_purchase_detail_id);

CREATE INDEX IF NOT EXISTS idx_course_purchase_revenue_shares_purchase_id
  ON course_purchase_revenue_shares(course_purchase_id);

CREATE INDEX IF NOT EXISTS idx_course_purchase_revenue_shares_instructor_user_id
  ON course_purchase_revenue_shares(instructor_user_id);

CREATE INDEX IF NOT EXISTS idx_course_purchase_revenue_shares_purchased_at
  ON course_purchase_revenue_shares(purchased_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_course_purchase_revenue_shares_purchased_at;
DROP INDEX IF EXISTS idx_course_purchase_revenue_shares_instructor_user_id;
DROP INDEX IF EXISTS idx_course_purchase_revenue_shares_purchase_id;
DROP INDEX IF EXISTS uq_course_purchase_revenue_shares_detail_id;
DROP TABLE IF EXISTS course_purchase_revenue_shares;
-- +goose StatementEnd
