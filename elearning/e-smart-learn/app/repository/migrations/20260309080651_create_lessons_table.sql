-- +goose Up
-- +goose StatementBegin
CREATE TABLE lessons (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  chapter_id UUID NOT NULL REFERENCES chapters(id) ON DELETE CASCADE,
  title VARCHAR(255) NOT NULL,
  content TEXT,
  video_url TEXT,
  is_able_to_preview BOOLEAN DEFAULT false,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

CREATE INDEX idx_lessons_chapter_id ON lessons(chapter_id);
CREATE INDEX idx_lessons_created_at ON lessons(created_at);
CREATE INDEX idx_lessons_deleted_at ON lessons(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_lessons_deleted_at;
DROP INDEX IF EXISTS idx_lessons_created_at;
DROP INDEX IF EXISTS idx_lessons_chapter_id;
DROP TABLE IF EXISTS lessons;
-- +goose StatementEnd
