-- +goose Up
-- +goose StatementBegin

-- Drop description column from chapters table
ALTER TABLE chapters DROP COLUMN IF EXISTS description;
ALTER TABLE chapters ADD COLUMN IF NOT EXISTS "order" INTEGER DEFAULT 0;

-- Create lesson_type enum if it doesn't exist
DO $$ BEGIN
  CREATE TYPE lesson_type AS ENUM ('video', 'document');
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

-- Add new columns to lessons table
ALTER TABLE lessons ADD COLUMN IF NOT EXISTS duration INTEGER DEFAULT 0;
ALTER TABLE lessons ADD COLUMN IF NOT EXISTS document_url TEXT;
ALTER TABLE lessons ADD COLUMN IF NOT EXISTS "order" INTEGER DEFAULT 0;
ALTER TABLE lessons ADD COLUMN IF NOT EXISTS "type" lesson_type DEFAULT 'video';

-- Drop content column from lessons table
ALTER TABLE lessons DROP COLUMN IF EXISTS content;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Revert changes to chapters table
ALTER TABLE chapters ADD COLUMN IF NOT EXISTS description TEXT;
ALTER TABLE chapters DROP COLUMN IF EXISTS "order";

-- Remove columns from lessons table and add content back
ALTER TABLE lessons DROP COLUMN IF EXISTS duration;
ALTER TABLE lessons DROP COLUMN IF EXISTS document_url;
ALTER TABLE lessons DROP COLUMN IF EXISTS "order";
ALTER TABLE lessons DROP COLUMN IF EXISTS "type";
ALTER TABLE lessons ADD COLUMN IF NOT EXISTS content TEXT;

-- Drop lesson_type enum
DROP TYPE IF EXISTS lesson_type CASCADE;

-- +goose StatementEnd
