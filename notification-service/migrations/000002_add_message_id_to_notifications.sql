-- +goose Up
ALTER TABLE notifications ADD COLUMN message_id VARCHAR(255) NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE notifications DROP COLUMN IF EXISTS message_id;
