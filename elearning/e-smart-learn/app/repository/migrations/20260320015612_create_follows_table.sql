-- +goose Up
-- +goose StatementBegin
CREATE TABLE follows (
	follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	followee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (follower_id, followee_id),
	CONSTRAINT chk_follows_not_self CHECK (follower_id <> followee_id)
);

CREATE INDEX idx_follows_follower_id ON follows(follower_id);
CREATE INDEX idx_follows_followee_id ON follows(followee_id);
CREATE INDEX idx_follows_created_at ON follows(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_follows_created_at;
DROP INDEX IF EXISTS idx_follows_followee_id;
DROP INDEX IF EXISTS idx_follows_follower_id;
DROP TABLE IF EXISTS follows;
-- +goose StatementEnd