package notification

import (
	"context"
	"database/sql"
	"errors"

	"github.com/anomalyco/hookah-store/notification-service/internal/errs"
	"github.com/anomalyco/hookah-store/notification-service/internal/models"
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{db: db}
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
