-- +goose Up
ALTER TABLE outbox ADD COLUMN type VARCHAR(50) NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE outbox DROP COLUMN IF EXISTS type;
