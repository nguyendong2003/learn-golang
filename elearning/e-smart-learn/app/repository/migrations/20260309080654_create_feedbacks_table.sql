-- +goose Up
-- +goose StatementBegin
CREATE TABLE feedbacks (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
  rate INT DEFAULT 5,
  content TEXT,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_feedbacks_user_id ON feedbacks(user_id);
CREATE INDEX idx_feedbacks_course_id ON feedbacks(course_id);
CREATE INDEX idx_feedbacks_created_at ON feedbacks(created_at);
CREATE INDEX idx_feedbacks_deleted_at ON feedbacks(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_feedbacks_deleted_at;
DROP INDEX IF EXISTS idx_feedbacks_created_at;
DROP INDEX IF EXISTS idx_feedbacks_course_id;
DROP INDEX IF EXISTS idx_feedbacks_user_id;
DROP TABLE IF EXISTS feedbacks;
-- +goose StatementEnd
