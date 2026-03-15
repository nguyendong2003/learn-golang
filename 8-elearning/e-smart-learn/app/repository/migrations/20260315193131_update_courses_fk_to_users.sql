-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.courses
DROP CONSTRAINT fk_courses_instructor_id;

ALTER TABLE public.courses
DROP COLUMN instructor_id;

ALTER TABLE public.courses
ADD COLUMN user_id uuid NOT NULL;

ALTER TABLE public.courses
ADD CONSTRAINT fk_courses_user_id
FOREIGN KEY (user_id)
REFERENCES public.users(id)
ON DELETE CASCADE;

CREATE INDEX idx_courses_user_id
ON public.courses(user_id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.courses
DROP CONSTRAINT fk_courses_user_id;

DROP INDEX IF EXISTS idx_courses_user_id;

ALTER TABLE public.courses
DROP COLUMN user_id;

ALTER TABLE public.courses
ADD COLUMN instructor_id uuid NOT NULL;

ALTER TABLE public.courses
ADD CONSTRAINT fk_courses_instructor_id
FOREIGN KEY (instructor_id)
REFERENCES public.instructor_profiles(id)
ON DELETE CASCADE;

-- +goose StatementEnd
