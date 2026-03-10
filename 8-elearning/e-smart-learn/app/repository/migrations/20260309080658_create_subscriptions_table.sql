-- +goose Up
-- +goose StatementBegin
CREATE TABLE subscriptions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
  billing_cycle VARCHAR(50) NOT NULL,
  status VARCHAR(50) DEFAULT 'active',
  started_at TIMESTAMP NOT NULL,
  ended_at TIMESTAMP,
  canceled_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_plan_id ON subscriptions(plan_id);
CREATE INDEX idx_subscriptions_created_at ON subscriptions(created_at);
CREATE INDEX idx_subscriptions_deleted_at ON subscriptions(deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_subscriptions_deleted_at;
DROP INDEX IF EXISTS idx_subscriptions_created_at;
DROP INDEX IF EXISTS idx_subscriptions_plan_id;
DROP INDEX IF EXISTS idx_subscriptions_user_id;
DROP TABLE IF EXISTS subscriptions;
-- +goose StatementEnd
