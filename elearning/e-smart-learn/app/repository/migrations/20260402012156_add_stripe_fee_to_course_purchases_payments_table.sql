-- +goose Up
-- +goose StatementBegin
ALTER TABLE course_purchases
  ADD COLUMN IF NOT EXISTS stripe_fee BIGINT NOT NULL DEFAULT 0;

ALTER TABLE payments
  ADD COLUMN IF NOT EXISTS stripe_fee BIGINT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE payments
  DROP COLUMN IF EXISTS stripe_fee;

ALTER TABLE course_purchases
  DROP COLUMN IF EXISTS stripe_fee;
-- +goose StatementEnd
