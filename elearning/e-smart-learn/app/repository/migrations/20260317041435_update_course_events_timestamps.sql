-- +goose Up
-- +goose StatementBegin
ALTER TABLE course_events 
  ALTER COLUMN start_time TYPE TIMESTAMPTZ,
  ALTER COLUMN end_time TYPE TIMESTAMPTZ;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE course_events 
  ALTER COLUMN start_time TYPE TIMESTAMP,
  ALTER COLUMN end_time TYPE TIMESTAMP;
-- +goose StatementEnd