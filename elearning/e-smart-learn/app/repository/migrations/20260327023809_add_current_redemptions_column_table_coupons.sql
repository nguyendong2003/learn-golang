-- +goose Up
-- +goose StatementBegin
ALTER TABLE coupons
  ADD COLUMN IF NOT EXISTS current_redemptions BIGINT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE coupons
  DROP COLUMN IF EXISTS current_redemptions;
-- +goose StatementEnd
