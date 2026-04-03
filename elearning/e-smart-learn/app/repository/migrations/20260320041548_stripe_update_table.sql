-- +goose Up
-- +goose StatementBegin

-- ======================
-- ENUMS
-- ======================
CREATE TYPE billing_cycle_enum AS ENUM ('monthly', 'yearly');

CREATE TYPE subscription_status_enum AS ENUM (
  'incomplete',
  'trialing',
  'active',
  'past_due',
  'canceled',
  'unpaid',
  'incomplete_expired'
);

CREATE TYPE payment_status_enum AS ENUM (
  'pending',
  'succeeded',
  'failed',
  'refunded'
);

-- ======================
-- USERS
-- ======================
ALTER TABLE users
  ADD COLUMN stripe_customer_id VARCHAR(255);

CREATE UNIQUE INDEX idx_users_stripe_customer_id 
  ON users(stripe_customer_id);


-- ======================
-- PLANS
-- ======================
ALTER TABLE plans
  ADD COLUMN billing_cycle billing_cycle_enum NOT NULL DEFAULT 'monthly',
  ADD COLUMN price DECIMAL(10,2) NOT NULL DEFAULT 0,
  ADD COLUMN stripe_price_id VARCHAR(255) NOT NULL,
  ADD COLUMN currency VARCHAR(10) DEFAULT 'usd';

CREATE UNIQUE INDEX idx_plans_stripe_price_id 
  ON plans(stripe_price_id);

CREATE INDEX idx_plans_billing_cycle 
  ON plans(billing_cycle);

-- drop legacy columns
ALTER TABLE plans
  DROP COLUMN IF EXISTS monthly_price,
  DROP COLUMN IF EXISTS yearly_price,
  DROP COLUMN IF EXISTS stripe_price_monthly_id,
  DROP COLUMN IF EXISTS stripe_price_yearly_id,
  DROP COLUMN IF EXISTS is_default;


-- ======================
-- SUBSCRIPTIONS
-- ======================
ALTER TABLE subscriptions
  ADD COLUMN stripe_subscription_id VARCHAR(255),
  ADD COLUMN stripe_customer_id VARCHAR(255),
  ADD COLUMN current_period_start TIMESTAMP,
  ADD COLUMN current_period_end TIMESTAMP,
  ADD COLUMN cancel_at_period_end BOOLEAN DEFAULT false;

-- convert billing_cycle -> enum
ALTER TABLE subscriptions
  ALTER COLUMN billing_cycle TYPE billing_cycle_enum
  USING billing_cycle::billing_cycle_enum;

-- FIX: drop default trước khi đổi enum
ALTER TABLE subscriptions
  ALTER COLUMN status DROP DEFAULT,
  ALTER COLUMN status TYPE subscription_status_enum
  USING status::subscription_status_enum,
  ALTER COLUMN status SET DEFAULT 'incomplete';

CREATE UNIQUE INDEX idx_subscriptions_stripe_subscription_id 
  ON subscriptions(stripe_subscription_id);

CREATE INDEX idx_subscriptions_stripe_customer_id 
  ON subscriptions(stripe_customer_id);


-- ======================
-- PAYMENTS
-- ======================
ALTER TABLE payments
  ADD COLUMN stripe_invoice_id VARCHAR(255),
  ADD COLUMN stripe_payment_intent VARCHAR(255),
  ADD COLUMN amount BIGINT DEFAULT 0,
  ADD COLUMN currency VARCHAR(10) DEFAULT 'usd',
  ADD COLUMN failure_reason TEXT,
  ADD COLUMN attempt_count BIGINT DEFAULT 0;

-- FIX: drop default trước khi đổi enum
ALTER TABLE payments
  ALTER COLUMN status DROP DEFAULT,
  ALTER COLUMN status TYPE payment_status_enum
  USING status::payment_status_enum,
  ALTER COLUMN status SET DEFAULT 'pending';

CREATE INDEX idx_payments_stripe_invoice_id 
  ON payments(stripe_invoice_id);

CREATE INDEX idx_payments_stripe_payment_intent 
  ON payments(stripe_payment_intent);


-- ======================
-- STRIPE EVENTS
-- ======================
CREATE TABLE stripe_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  event_id VARCHAR(255) UNIQUE NOT NULL,
  event_type VARCHAR(100) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

CREATE INDEX idx_stripe_events_event_type 
  ON stripe_events(event_type);

CREATE INDEX idx_stripe_events_deleted_at 
  ON stripe_events(deleted_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- STRIPE EVENTS
DROP INDEX IF EXISTS idx_stripe_events_deleted_at;
DROP INDEX IF EXISTS idx_stripe_events_event_type;
DROP TABLE IF EXISTS stripe_events;

-- PAYMENTS
DROP INDEX IF EXISTS idx_payments_stripe_payment_intent;
DROP INDEX IF EXISTS idx_payments_stripe_invoice_id;

ALTER TABLE payments
  ALTER COLUMN status TYPE VARCHAR(50) USING status::text,
  ALTER COLUMN status SET DEFAULT 'pending',
  DROP COLUMN IF EXISTS attempt_count,
  DROP COLUMN IF EXISTS failure_reason,
  DROP COLUMN IF EXISTS currency,
  DROP COLUMN IF EXISTS amount,
  DROP COLUMN IF EXISTS stripe_payment_intent,
  DROP COLUMN IF EXISTS stripe_invoice_id;

-- SUBSCRIPTIONS
DROP INDEX IF EXISTS idx_subscriptions_stripe_customer_id;
DROP INDEX IF EXISTS idx_subscriptions_stripe_subscription_id;

ALTER TABLE subscriptions
  ALTER COLUMN status TYPE VARCHAR(50) USING status::text,
  ALTER COLUMN status SET DEFAULT 'active',
  ALTER COLUMN billing_cycle TYPE VARCHAR(50) USING billing_cycle::text,
  DROP COLUMN IF EXISTS cancel_at_period_end,
  DROP COLUMN IF EXISTS current_period_end,
  DROP COLUMN IF EXISTS current_period_start,
  DROP COLUMN IF EXISTS stripe_customer_id,
  DROP COLUMN IF EXISTS stripe_subscription_id;

-- PLANS
DROP INDEX IF EXISTS idx_plans_billing_cycle;
DROP INDEX IF EXISTS idx_plans_stripe_price_id;

ALTER TABLE plans
  ALTER COLUMN billing_cycle TYPE VARCHAR(20) USING billing_cycle::text,
  ADD COLUMN IF NOT EXISTS monthly_price DECIMAL(10,2) DEFAULT 0,
  ADD COLUMN IF NOT EXISTS yearly_price DECIMAL(10,2) DEFAULT 0,
  ADD COLUMN IF NOT EXISTS stripe_price_monthly_id VARCHAR(255),
  ADD COLUMN IF NOT EXISTS stripe_price_yearly_id VARCHAR(255),
  ADD COLUMN IF NOT EXISTS is_default BOOLEAN DEFAULT false;

ALTER TABLE plans
  DROP COLUMN IF EXISTS stripe_price_id,
  DROP COLUMN IF EXISTS price,
  DROP COLUMN IF EXISTS billing_cycle,
  DROP COLUMN IF EXISTS currency;

-- USERS
DROP INDEX IF EXISTS idx_users_stripe_customer_id;

ALTER TABLE users 
  DROP COLUMN IF EXISTS stripe_customer_id;

-- ENUMS
DROP TYPE IF EXISTS payment_status_enum;
DROP TYPE IF EXISTS subscription_status_enum;
DROP TYPE IF EXISTS billing_cycle_enum;

-- +goose StatementEnd