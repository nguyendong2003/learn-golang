-- +goose Up
-- +goose StatementBegin
ALTER TABLE courses DROP CONSTRAINT courses_teacher_id_fkey;
DROP INDEX idx_courses_teacher_id;
ALTER TABLE courses RENAME COLUMN teacher_id TO instructor_id;
ALTER TABLE courses ADD CONSTRAINT fk_courses_instructor_id 
  FOREIGN KEY (instructor_id) REFERENCES instructor_profiles(id) ON DELETE CASCADE;

CREATE INDEX idx_courses_instructor_id ON courses(instructor_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE courses DROP CONSTRAINT fk_courses_instructor_id;
DROP INDEX idx_courses_instructor_id;
ALTER TABLE courses RENAME COLUMN instructor_id TO teacher_id;

ALTER TABLE courses ADD CONSTRAINT courses_teacher_id_fkey 
  FOREIGN KEY (teacher_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX idx_courses_teacher_id ON courses(teacher_id);
-- +goose StatementEnd
