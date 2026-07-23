package repository

import (
	"context"

	"github.com/anomalyco/hookah-store/notification-service/internal/models"
	"github.com/google/uuid"
)

type NotificationRepository interface {
	Create(ctx context.Context, message *models.Notification) (id uuid.UUID, err error)
	UpdateMessageID(ctx context.Context, id uuid.UUID, msgID string) error
	UpdateStatus(ctx context.Context, email string, msgID string, status models.MsgStatus) error
	GetByUserID(ctx context.Context, userID string) ([]models.Notification, error)
	GetByEmail(ctx context.Context, email string) ([]models.Notification, error)
}
