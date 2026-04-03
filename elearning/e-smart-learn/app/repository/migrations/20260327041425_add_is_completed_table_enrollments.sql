-- +goose Up
ALTER TABLE enrollments
ADD COLUMN is_completed boolean NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE enrollments
DROP COLUMN IF EXISTS is_completed;
