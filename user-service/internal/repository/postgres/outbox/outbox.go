package outbox

import (
	"context"
	"fmt"

	"github.com/anomalyco/hookah-store/user-service/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

func (o *Repo) SaveEvent(ctx context.Context, tx *sqlx.Tx, event *models.OutboxEvent) error {
	query := `
		INSERT INTO outbox (id, topic, key, type, payload, created_at, published)
		VALUES (:id, :topic, :key, :type, :payload, :created_at, :published)`

	_, err := tx.NamedExecContext(ctx, query, event)
	if err != nil {
		return fmt.Errorf("failed save outbox event: %w", err)
	}

	return nil
}

func (o *Repo) FetchUnpublishedEvents(limit int) ([]models.OutboxEvent, error) {
	query := `SELECT id, topic, key, type, payload, created_at
			FROM outbox
			WHERE published = false
			ORDER BY created_at ASC
			LIMIT $1`

	rows, err := o.db.Queryx(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch unpublished events: %w", err)
	}
	defer rows.Close()

	var events []models.OutboxEvent

	for rows.Next() {
		var event models.OutboxEvent
		if err := rows.StructScan(&event); err != nil {
			return nil, fmt.Errorf("failed to scan outbox event: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return events, nil
}

func (o *Repo) MarkPublishedEvents(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE outbox SET published = true WHERE id = $1`

	result, err := o.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("mark event published: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("event %s not found", id)
	}

	return nil
}
