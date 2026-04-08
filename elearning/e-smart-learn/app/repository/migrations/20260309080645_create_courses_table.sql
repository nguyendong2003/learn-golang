-- +goose Up
-- +goose StatementBegin
CREATE TABLE courses (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title VARCHAR(255) NOT NULL,
  description TEXT,
  image TEXT,
  slug VARCHAR(255) NOT NULL UNIQUE,
  teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  duration INT DEFAULT 0,
  category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
  price DECIMAL(10, 2) DEFAULT 0,
  old_price DECIMAL(10, 2) DEFAULT 0,
  average_rate DECIMAL(3, 2) DEFAULT 0,
  status VARCHAR(50) DEFAULT 'draft',
  total_student BIGINT DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_courses_teacher_id ON courses(teacher_id);
CREATE INDEX idx_courses_category_id ON courses(category_id);
CREATE INDEX idx_courses_created_at ON courses(created_at);
CREATE INDEX idx_courses_deleted_at ON courses(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_courses_deleted_at;
DROP INDEX IF EXISTS idx_courses_created_at;
DROP INDEX IF EXISTS idx_courses_category_id;
DROP INDEX IF EXISTS idx_courses_teacher_id;
DROP TABLE IF EXISTS courses;
-- +goose StatementEnd
