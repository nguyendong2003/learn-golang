-- +goose Up
-- +goose StatementBegin
CREATE TABLE blogs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title VARCHAR(255) NOT NULL,
  content TEXT NOT NULL,
  slug VARCHAR(255) NOT NULL UNIQUE,
  author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  view_total BIGINT DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

CREATE INDEX idx_blogs_author_id ON blogs(author_id);
CREATE INDEX idx_blogs_created_at ON blogs(created_at);
CREATE INDEX idx_blogs_deleted_at ON blogs(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_blogs_deleted_at;
DROP INDEX IF EXISTS idx_blogs_created_at;
DROP INDEX IF EXISTS idx_blogs_author_id;
DROP TABLE IF EXISTS blogs;
-- +goose StatementEnd
