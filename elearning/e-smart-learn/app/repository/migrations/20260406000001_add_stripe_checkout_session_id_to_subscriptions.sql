-- +goose Up
-- +goose StatementBegin
ALTER TABLE subscriptions
  ADD COLUMN IF NOT EXISTS stripe_checkout_session_id VARCHAR(255);

CREATE INDEX IF NOT EXISTS idx_subscriptions_stripe_checkout_session_id
  ON subscriptions(stripe_checkout_session_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_subscriptions_stripe_checkout_session_id;

ALTER TABLE subscriptions
  DROP COLUMN IF EXISTS stripe_checkout_session_id;
-- +goose StatementEnd