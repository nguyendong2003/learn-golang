-- +goose Up
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS role;
ALTER TABLE users ADD COLUMN role_id UUID;

UPDATE users SET role_id = (SELECT id FROM roles LIMIT 1) WHERE role_id IS NULL;

ALTER TABLE users ALTER COLUMN role_id SET NOT NULL;

ALTER TABLE users ADD CONSTRAINT fk_users_role_id 
  FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT;

CREATE INDEX idx_users_role_id ON users(role_id);

DROP TYPE IF EXISTS user_role CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Create the type again
CREATE TYPE user_role AS ENUM ('student', 'teacher', 'admin');
ALTER TABLE users ADD COLUMN role user_role DEFAULT 'student';
DROP INDEX IF EXISTS idx_users_role_id;
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_role_id;
ALTER TABLE users DROP COLUMN role_id;
-- +goose StatementEnd
