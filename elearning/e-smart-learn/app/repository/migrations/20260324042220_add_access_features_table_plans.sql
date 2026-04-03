-- +goose Up
ALTER TABLE plans
ADD COLUMN IF NOT EXISTS access_features jsonb NOT NULL DEFAULT '[]'::jsonb;

-- create index to help queries filtering by features if needed
CREATE INDEX IF NOT EXISTS idx_plans_access_features ON plans USING gin (access_features);

-- +goose Down
DROP INDEX IF EXISTS idx_plans_access_features;
ALTER TABLE plans DROP COLUMN IF EXISTS access_features;
