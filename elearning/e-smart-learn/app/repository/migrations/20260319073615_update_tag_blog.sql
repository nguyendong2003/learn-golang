-- +goose Up
-- +goose StatementBegin
ALTER TABLE blogs ADD COLUMN IF NOT EXISTS tags jsonb;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE blogs DROP COLUMN IF EXISTS tags;
-- +goose StatementEnd
