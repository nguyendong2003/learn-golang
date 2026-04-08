-- +goose Up
-- +goose StatementBegin
ALTER TABLE courses
  ADD COLUMN IF NOT EXISTS stripe_product_id VARCHAR(255),
  ADD COLUMN IF NOT EXISTS stripe_price_id VARCHAR(255),
  ADD COLUMN IF NOT EXISTS stripe_currency VARCHAR(10) DEFAULT 'usd',
  ADD COLUMN IF NOT EXISTS stripe_amount BIGINT DEFAULT 0,
  ADD COLUMN IF NOT EXISTS stripe_synced_at TIMESTAMPTZ,
  DROP COLUMN IF EXISTS old_price;

CREATE INDEX IF NOT EXISTS idx_courses_stripe_product_id ON courses(stripe_product_id);
CREATE INDEX IF NOT EXISTS idx_courses_stripe_price_id ON courses(stripe_price_id);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_courses_stripe_price_id;
DROP INDEX IF EXISTS idx_courses_stripe_product_id;
ALTER TABLE courses
  DROP COLUMN IF EXISTS stripe_product_id,
  DROP COLUMN IF EXISTS stripe_price_id,
  DROP COLUMN IF EXISTS stripe_currency,
  DROP COLUMN IF EXISTS stripe_amount,
  DROP COLUMN IF EXISTS stripe_synced_at,
  ADD COLUMN IF NOT EXISTS old_price DECIMAL(10, 2) DEFAULT 0;
-- +goose StatementEnd
