-- +goose Up
-- +goose StatementBegin
ALTER TABLE course_events
	ADD COLUMN IF NOT EXISTS attendee_emails JSONB NOT NULL DEFAULT '[]'::jsonb;

ALTER TABLE course_events
	DROP COLUMN IF EXISTS email;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE course_events
	ADD COLUMN IF NOT EXISTS email VARCHAR(255);

ALTER TABLE course_events
	DROP COLUMN IF EXISTS attendee_emails;
-- +goose StatementEnd
