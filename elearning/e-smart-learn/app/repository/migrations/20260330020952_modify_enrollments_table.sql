-- +goose Up
-- +goose StatementBegin
ALTER TABLE enrollments
	ADD COLUMN id UUID NOT NULL DEFAULT gen_random_uuid();

ALTER TABLE enrollments
	DROP CONSTRAINT IF EXISTS enrollments_pkey;

ALTER TABLE enrollments
	ADD CONSTRAINT enrollments_pkey PRIMARY KEY (id);

ALTER TABLE enrollments
	ADD COLUMN canceled_at TIMESTAMP;

ALTER TABLE enrollments
    ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
    
ALTER TABLE enrollments
    ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

CREATE TYPE enrollment_type_enum AS ENUM ('course_purchase', 'subscription');

ALTER TABLE enrollments
    ADD COLUMN type enrollment_type_enum DEFAULT 'subscription';

CREATE INDEX idx_enrollments_canceled_at ON enrollments(canceled_at);
CREATE INDEX idx_enrollments_type ON enrollments(type);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE enrollments
	DROP CONSTRAINT IF EXISTS enrollments_pkey;

ALTER TABLE enrollments
	DROP COLUMN IF EXISTS id;

ALTER TABLE enrollments
	ADD CONSTRAINT enrollments_pkey PRIMARY KEY (user_id, course_id);

ALTER TABLE enrollments
	DROP COLUMN IF EXISTS canceled_at;

DROP INDEX IF EXISTS idx_enrollments_canceled_at;
DROP INDEX IF EXISTS idx_enrollments_type;

ALTER TABLE enrollments
	DROP COLUMN IF EXISTS type;
	
DROP TYPE IF EXISTS enrollment_type_enum;
-- +goose StatementEnd
