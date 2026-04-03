-- +goose Up
-- +goose StatementBegin
ALTER TABLE plans
  ADD COLUMN IF NOT EXISTS stripe_product_id VARCHAR(255) NOT NULL UNIQUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE plans
  DROP COLUMN IF EXISTS stripe_product_id;
-- +goose StatementEnd