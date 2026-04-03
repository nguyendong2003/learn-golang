-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN IF NOT EXISTS oauth_provider VARCHAR(50),
ADD COLUMN IF NOT EXISTS oauth_id VARCHAR(255);

CREATE UNIQUE INDEX IF NOT EXISTS ux_users_oauth_provider_id
ON users (oauth_provider, oauth_id)
WHERE oauth_provider IS NOT NULL AND oauth_id IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS ux_users_oauth_provider_id;
ALTER TABLE users
DROP COLUMN IF EXISTS oauth_id,
DROP COLUMN IF EXISTS oauth_provider;
-- +goose StatementEnd
