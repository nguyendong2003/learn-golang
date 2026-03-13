-- +goose Up
-- +goose StatementBegin
ALTER TABLE "courses" DROP CONSTRAINT IF EXISTS "fk_courses_instructor_id";
ALTER TABLE "courses" 
ADD CONSTRAINT "fk_courses_users" id
FOREIGN KEY ("instructor_id") 
REFERENCES "users"("id") 
ON DELETE SET NULL; -- Hoặc CASCADE tùy vào logic nghiệp vụ của bạn
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "courses" DROP CONSTRAINT IF EXISTS "fk_courses_users";
ALTER TABLE "courses" 
ADD CONSTRAINT "fk_courses_instructor_id" 
FOREIGN KEY ("instructor_id") 
REFERENCES "instructor_profiles"("id");
-- +goose StatementEnd