-- +goose Up
ALTER TABLE instructor_profiles
ADD COLUMN IF NOT EXISTS category_id uuid;

-- add foreign key constraint to categories(id)
ALTER TABLE instructor_profiles
ADD CONSTRAINT fk_instructor_profiles_category
FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL;

-- index for category_id
CREATE INDEX IF NOT EXISTS idx_instructor_profiles_category_id ON instructor_profiles(category_id);

-- +goose Down
ALTER TABLE instructor_profiles DROP CONSTRAINT fk_instructor_profiles_category;
DROP INDEX idx_instructor_profiles_category_id;
ALTER TABLE instructor_profiles DROP COLUMN IF EXISTS category_id;
