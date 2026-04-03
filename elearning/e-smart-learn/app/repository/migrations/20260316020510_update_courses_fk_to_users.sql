-- +goose Up
-- +goose StatementBegin

ALTER TABLE courses DROP CONSTRAINT IF EXISTS fk_courses_instructor_id;

ALTER TABLE courses DROP COLUMN IF EXISTS instructor_id;

ALTER TABLE courses ADD COLUMN user_id uuid NOT NULL;

ALTER TABLE courses
ADD CONSTRAINT fk_courses_user_id
FOREIGN KEY (user_id)
REFERENCES users(id)
ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_courses_user_id ON courses(user_id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE courses DROP CONSTRAINT IF EXISTS fk_courses_user_id;

DROP INDEX IF EXISTS idx_courses_user_id;

ALTER TABLE courses DROP COLUMN IF EXISTS user_id;

ALTER TABLE courses ADD COLUMN instructor_id uuid NOT NULL;

ALTER TABLE courses
ADD CONSTRAINT fk_courses_instructor_id
FOREIGN KEY (instructor_id)
REFERENCES instructor_profiles(id)
ON DELETE CASCADE;

-- +goose StatementEnd