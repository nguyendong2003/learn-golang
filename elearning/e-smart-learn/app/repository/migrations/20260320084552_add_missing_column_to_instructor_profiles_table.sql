-- +goose Up
-- +goose StatementBegin
ALTER TABLE instructor_profiles ADD COLUMN IF NOT EXISTS years_of_experience int DEFAULT 0;
ALTER TABLE instructor_profiles ADD COLUMN IF NOT EXISTS cv_url text;
ALTER TABLE instructor_profiles ADD COLUMN IF NOT EXISTS portfolio_url text;
ALTER TABLE instructor_profiles ADD COLUMN IF NOT EXISTS certifications jsonb;
ALTER TABLE instructor_profiles ADD COLUMN IF NOT EXISTS status text DEFAULT 'pending_review';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE instructor_profiles DROP COLUMN IF EXISTS years_of_experience;
ALTER TABLE instructor_profiles DROP COLUMN IF EXISTS cv_url;
ALTER TABLE instructor_profiles DROP COLUMN IF EXISTS portfolio_url;
ALTER TABLE instructor_profiles DROP COLUMN IF EXISTS certifications;
ALTER TABLE instructor_profiles DROP COLUMN IF EXISTS status;
-- +goose StatementEnd