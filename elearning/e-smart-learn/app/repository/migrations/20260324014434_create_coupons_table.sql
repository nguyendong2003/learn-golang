-- +goose Up
-- +goose StatementBegin
CREATE TYPE discount_type_enum AS ENUM ('percent', 'amount');

CREATE TABLE IF NOT EXISTS coupons (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code VARCHAR(100) NOT NULL UNIQUE,
  stripe_coupon_id VARCHAR(255) NOT NULL,
  stripe_promotion_code_id VARCHAR(255) NOT NULL UNIQUE,
  discount_type discount_type_enum NOT NULL,
  discount_value BIGINT DEFAULT 0,
  currency VARCHAR(10) DEFAULT 'usd',
  max_redemptions BIGINT DEFAULT 0,
  is_active BOOLEAN DEFAULT true,
  expires_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_coupons_deleted_at ON coupons(deleted_at);
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_coupons_deleted_at;
DROP TABLE IF EXISTS coupons;
DROP TYPE IF EXISTS discount_type_enum;
-- +goose StatementEnd
