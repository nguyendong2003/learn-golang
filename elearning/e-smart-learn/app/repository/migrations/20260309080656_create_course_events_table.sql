-- +goose Up
-- +goose StatementBegin
CREATE TABLE course_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  meeting_url TEXT,
  start_time TIMESTAMPTZ NOT NULL,
  end_time TIMESTAMPTZ NOT NULL,
  notification_before_start INT DEFAULT 15,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_course_events_course_id ON course_events(course_id);
CREATE INDEX idx_course_events_start_time ON course_events(start_time);
CREATE INDEX idx_course_events_created_at ON course_events(created_at);
CREATE INDEX idx_course_events_deleted_at ON course_events(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_course_events_deleted_at;
DROP INDEX IF EXISTS idx_course_events_created_at;
DROP INDEX IF EXISTS idx_course_events_start_time;
DROP INDEX IF EXISTS idx_course_events_course_id;
DROP TABLE IF EXISTS course_events;
-- +goose StatementEnd
