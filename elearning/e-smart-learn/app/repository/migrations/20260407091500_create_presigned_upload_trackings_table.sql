-- +goose Up
CREATE TABLE presigned_upload_trackings (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    object_url TEXT NOT NULL UNIQUE,
    object_name TEXT NOT NULL,
    filetype VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_presign_trackings_status_created ON presigned_upload_trackings(status, created_at);

-- +goose Down
DROP TABLE IF EXISTS presigned_upload_trackings;
