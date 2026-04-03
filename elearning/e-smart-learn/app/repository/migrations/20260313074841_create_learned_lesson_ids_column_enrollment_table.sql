-- +goose Up
-- Add a uuid[] column to store learned lesson IDs for a user's enrollment
ALTER TABLE enrollments
	ADD COLUMN IF NOT EXISTS learned_lesson_ids uuid[] DEFAULT '{}'::uuid[];

-- +goose Down
-- Remove the learned_lesson_ids column
ALTER TABLE enrollments
	DROP COLUMN IF EXISTS learned_lesson_ids;
