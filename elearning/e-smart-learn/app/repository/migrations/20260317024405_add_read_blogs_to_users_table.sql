-- +goose Up
-- Add a uuid[] column to store read blog IDs for a user
ALTER TABLE users
	ADD COLUMN IF NOT EXISTS read_blog_ids uuid[] DEFAULT '{}'::uuid[];

-- +goose Down
-- Remove the read_blog_ids column
ALTER TABLE users
	DROP COLUMN IF EXISTS read_blog_ids;
