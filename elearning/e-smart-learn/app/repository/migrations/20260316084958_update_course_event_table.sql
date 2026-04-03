-- +goose Up
-- +goose StatementBegin
ALTER TABLE course_events
DROP COLUMN meeting_url;

ALTER TABLE course_events
ADD COLUMN location VARCHAR(255);

ALTER TABLE course_events
ADD COLUMN room_token VARCHAR(255) UNIQUE;

ALTER TABLE course_events
ADD COLUMN email VARCHAR(255);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE course_events
DROP COLUMN location;

ALTER TABLE course_events
DROP COLUMN room_token;

ALTER TABLE course_events
DROP COLUMN email;

ALTER TABLE course_events
ADD COLUMN meeting_url TEXT;
-- +goose StatementEnd