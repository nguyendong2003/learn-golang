-- +goose Up
-- +goose StatementBegin
ALTER TABLE blogs ADD COLUMN IF NOT EXISTS image_url TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE blogs DROP COLUMN IF EXISTS image_url;
-- +goose StatementEnd
