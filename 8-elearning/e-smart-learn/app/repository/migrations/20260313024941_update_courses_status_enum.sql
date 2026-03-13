-- +goose Up
-- +goose StatementBegin

CREATE TYPE course_status AS ENUM (
  'draft',
  'pending_review',
  'published',
  'rejected',
  'archived'
);

ALTER TABLE courses
  ALTER COLUMN status DROP DEFAULT;

ALTER TABLE courses
  ALTER COLUMN status TYPE course_status
  USING status::course_status;

ALTER TABLE courses
  ALTER COLUMN status SET DEFAULT 'draft',
  ALTER COLUMN status SET NOT NULL;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE courses
  ALTER COLUMN status DROP DEFAULT,
  ALTER COLUMN status DROP NOT NULL;

ALTER TABLE courses
  ALTER COLUMN status TYPE VARCHAR(50)
  USING status::text;

ALTER TABLE courses
  ALTER COLUMN status SET DEFAULT 'draft';

DROP TYPE IF EXISTS course_status;

-- +goose StatementEnd