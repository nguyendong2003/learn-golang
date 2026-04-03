-- +goose Up
-- +goose StatementBegin
ALTER TABLE blogs ADD COLUMN status varchar(20) NOT NULL DEFAULT 'draft';
ALTER TABLE blogs ADD COLUMN scheduled_at timestamp with time zone;
ALTER TABLE blogs ADD COLUMN published_at timestamp with time zone;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE blogs DROP COLUMN IF EXISTS published_at;
ALTER TABLE blogs DROP COLUMN IF EXISTS scheduled_at;
ALTER TABLE blogs DROP COLUMN IF EXISTS status;
-- +goose StatementEnd
