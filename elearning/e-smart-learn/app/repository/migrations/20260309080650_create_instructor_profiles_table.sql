-- +goose Up
-- +goose StatementBegin
CREATE TABLE instructor_profiles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
  bio TEXT,
  education TEXT,
  rating_avg DECIMAL(3, 2) DEFAULT 0,
  total_student BIGINT DEFAULT 0,
  total_course BIGINT DEFAULT 0,
  balance DECIMAL(15, 2) DEFAULT 0,
  linkedin_url VARCHAR(255),
  youtube_url VARCHAR(255),
  instagram_url VARCHAR(255),
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_instructor_profiles_user_id ON instructor_profiles(user_id);
CREATE INDEX idx_instructor_profiles_created_at ON instructor_profiles(created_at);
CREATE INDEX idx_instructor_profiles_deleted_at ON instructor_profiles(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_instructor_profiles_deleted_at;
DROP INDEX IF EXISTS idx_instructor_profiles_created_at;
DROP INDEX IF EXISTS idx_instructor_profiles_user_id;
DROP TABLE IF EXISTS instructor_profiles;
-- +goose StatementEnd
