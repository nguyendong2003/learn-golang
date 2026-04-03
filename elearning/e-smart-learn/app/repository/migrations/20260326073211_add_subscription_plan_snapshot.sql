-- +goose Up
-- +goose StatementBegin
ALTER TABLE subscriptions
  ADD COLUMN IF NOT EXISTS plan_name VARCHAR(100) NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS plan_description TEXT,
  ADD COLUMN IF NOT EXISTS plan_price DECIMAL(10,2) NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS plan_currency VARCHAR(10) NOT NULL DEFAULT 'usd',
  ADD COLUMN IF NOT EXISTS plan_stripe_price_id VARCHAR(255);

CREATE INDEX IF NOT EXISTS idx_subscriptions_plan_stripe_price_id ON subscriptions(plan_stripe_price_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_subscriptions_plan_stripe_price_id;

ALTER TABLE subscriptions
  DROP COLUMN IF EXISTS plan_stripe_price_id,
  DROP COLUMN IF EXISTS plan_currency,
  DROP COLUMN IF EXISTS plan_price,
  DROP COLUMN IF EXISTS plan_description,
  DROP COLUMN IF EXISTS plan_name;
-- +goose StatementEnd
