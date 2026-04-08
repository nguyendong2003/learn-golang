-- +goose Up
-- +goose StatementBegin
ALTER TABLE payments
  ADD COLUMN IF NOT EXISTS billing_period_start TIMESTAMP,
  ADD COLUMN IF NOT EXISTS billing_period_end TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_payments_billing_period_start ON payments(billing_period_start);
CREATE INDEX IF NOT EXISTS idx_payments_billing_period_end ON payments(billing_period_end);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_payments_billing_period_end;
DROP INDEX IF EXISTS idx_payments_billing_period_start;

ALTER TABLE payments
  DROP COLUMN IF EXISTS billing_period_end,
  DROP COLUMN IF EXISTS billing_period_start;
-- +goose StatementEnd
