-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS carts (
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  
  PRIMARY KEY (user_id, course_id)
);

CREATE INDEX IF NOT EXISTS idx_carts_user_id ON carts(user_id);
CREATE INDEX IF NOT EXISTS idx_cart_course_id ON carts(course_id);
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_carts_course_id;
DROP INDEX IF EXISTS idx_carts_user_id;
DROP TABLE IF EXISTS carts;
-- +goose StatementEnd
