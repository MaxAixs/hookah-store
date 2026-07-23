package notification

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/anomalyco/hookah-store/notification-service/internal/errs"
	"github.com/anomalyco/hookah-store/notification-service/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, notification *models.Notification) (uuid.UUID, error) {
	query := `
		INSERT INTO notifications (id, user_id, event_id, email, subject, status, created_at, updated_at)
		VALUES (:id, :user_id, :event_id, :email, :subject, :status, :created_at, :updated_at)
		RETURNING id`

	rows, err := r.db.NamedQueryContext(ctx, query, notification)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create notification: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&notification.ID); err != nil {
			return uuid.Nil, fmt.Errorf("failed to scan returned id: %w", err)
		}
	}

	return notification.ID, nil
}

func (r *Repo) UpdateMessageID(ctx context.Context, id uuid.UUID, msgID string) error {
	query := `UPDATE notifications SET message_id = $1 WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, msgID, id)
	if err != nil {
		return fmt.Errorf("failed to update message_id: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errs.ErrUserNotFound
	}

	return nil
}

func (r *Repo) UpdateStatus(ctx context.Context, email string, msgID string, status models.MsgStatus) error {
	query := `UPDATE notifications SET status = $1, updated_at = NOW() WHERE message_id = $2 AND email = $3`

	result, err := r.db.ExecContext(ctx, query, string(status), msgID, email)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("notification with message_id %s and email %s not found", msgID, email)
	}

	return nil
}

func (r *Repo) GetByUserID(ctx context.Context, userID string) ([]models.Notification, error) {
	query := `
		SELECT id, user_id, email, subject, status, created_at, updated_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC`

	var notifications []models.Notification
	err := r.db.SelectContext(ctx, &notifications, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
	}

	return notifications, nil
}

func (r *Repo) GetByEmail(ctx context.Context, email string) ([]models.Notification, error) {
	query := `
		SELECT id, user_id, email, subject, status, created_at, updated_at
		FROM notifications
		WHERE email = $1
		ORDER BY created_at DESC`

	var notifications []models.Notification
	err := r.db.SelectContext(ctx, &notifications, query, email)
	if err != nil {
		return nil, errs.ErrEmailNotFound
	}

	return notifications, nil
}
