-- +goose Up
CREATE TABLE outbox (
    id         UUID PRIMARY KEY,
    topic      VARCHAR(255) NOT NULL,
    key        VARCHAR(255) NOT NULL,
    payload    JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    published  BOOLEAN NOT NULL DEFAULT false
);

CREATE INDEX idx_outbox_unpublished ON outbox (created_at) WHERE published = false;

-- +goose Down
DROP TABLE IF EXISTS outbox;
