-- +goose Up
-- +goose StatementBegin
ALTER TABLE blogs
	ADD COLUMN IF NOT EXISTS category_id UUID;

ALTER TABLE blogs
	ADD CONSTRAINT fk_blogs_category_id 
	FOREIGN KEY (category_id) 
	REFERENCES categories(id) 
	ON DELETE SET NULL;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE blogs DROP CONSTRAINT IF EXISTS fk_blogs_category_id;
ALTER TABLE blogs DROP COLUMN IF EXISTS category_id;
-- +goose StatementEnd