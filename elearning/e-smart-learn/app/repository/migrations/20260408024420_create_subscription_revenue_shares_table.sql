-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS subscription_revenue_shares (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  payment_id UUID NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
  subscription_id UUID NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
  subscriber_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  instructor_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  total_active_courses INTEGER NOT NULL DEFAULT 0,
  instructor_active_courses INTEGER NOT NULL DEFAULT 0,
  payment_amount BIGINT NOT NULL DEFAULT 0,
  payment_stripe_fee BIGINT NOT NULL DEFAULT 0,
  allocated_amount BIGINT NOT NULL DEFAULT 0,
  instructor_gross BIGINT NOT NULL DEFAULT 0,
  platform_gross BIGINT NOT NULL DEFAULT 0,
  allocated_stripe_fee BIGINT NOT NULL DEFAULT 0,
  instructor_net BIGINT NOT NULL DEFAULT 0,
  platform_net BIGINT NOT NULL DEFAULT 0,
  billing_period_start TIMESTAMPTZ,
  billing_period_end TIMESTAMPTZ,
  paid_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_subscription_revenue_shares_payment_instructor
  ON subscription_revenue_shares(payment_id, instructor_user_id);

CREATE INDEX IF NOT EXISTS idx_subscription_revenue_shares_payment_id
  ON subscription_revenue_shares(payment_id);

CREATE INDEX IF NOT EXISTS idx_subscription_revenue_shares_instructor_user_id
  ON subscription_revenue_shares(instructor_user_id);

CREATE INDEX IF NOT EXISTS idx_subscription_revenue_shares_subscription_id
  ON subscription_revenue_shares(subscription_id);

CREATE INDEX IF NOT EXISTS idx_subscription_revenue_shares_paid_at
  ON subscription_revenue_shares(paid_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_subscription_revenue_shares_paid_at;
DROP INDEX IF EXISTS idx_subscription_revenue_shares_subscription_id;
DROP INDEX IF EXISTS idx_subscription_revenue_shares_instructor_user_id;
DROP INDEX IF EXISTS idx_subscription_revenue_shares_payment_id;
DROP INDEX IF EXISTS uq_subscription_revenue_shares_payment_instructor;
DROP TABLE IF EXISTS subscription_revenue_shares;
-- +goose StatementEnd
