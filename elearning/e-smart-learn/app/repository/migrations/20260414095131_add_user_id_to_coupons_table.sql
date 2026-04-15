-- +goose Up
-- +goose StatementBegin
ALTER TABLE coupons
    ADD COLUMN IF NOT EXISTS user_id uuid;

CREATE INDEX IF NOT EXISTS idx_coupons_user_id ON coupons(user_id);
ALTER TABLE coupons
    ADD CONSTRAINT fk_coupons_user
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE coupons DROP CONSTRAINT IF EXISTS fk_coupons_user;
DROP INDEX IF EXISTS idx_coupons_user_id;
ALTER TABLE coupons DROP COLUMN IF EXISTS user_id;
-- +goose StatementEnd
