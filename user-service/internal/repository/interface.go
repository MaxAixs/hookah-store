package repository

import (
	"context"

	"github.com/anomalyco/hookah-store/user-service/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(ctx context.Context, tx *sqlx.Tx, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	UpdatePassword(ctx context.Context, tx *sqlx.Tx, email string, newPassword string) error
}

type OutBoxRepository interface {
	SaveEvent(ctx context.Context, tx *sqlx.Tx, event *models.OutboxEvent) error
	FetchUnpublishedEvents(limit int) ([]models.OutboxEvent, error)
	MarkPublishedEvents(ctx context.Context, id uuid.UUID) error
}
