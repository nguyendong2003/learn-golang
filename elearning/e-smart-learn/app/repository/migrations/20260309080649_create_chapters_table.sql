-- +goose Up
-- +goose StatementBegin
CREATE TABLE chapters (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_chapters_course_id ON chapters(course_id);
CREATE INDEX idx_chapters_created_at ON chapters(created_at);
CREATE INDEX idx_chapters_deleted_at ON chapters(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_chapters_deleted_at;
DROP INDEX IF EXISTS idx_chapters_created_at;
DROP INDEX IF EXISTS idx_chapters_course_id;
DROP TABLE IF EXISTS chapters;
-- +goose StatementEnd
