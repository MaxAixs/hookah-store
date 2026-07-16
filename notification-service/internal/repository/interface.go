package repository

import (
	"context"

	"github.com/anomalyco/hookah-store/notification-service/internal/models"
)

type NotificationRepository interface {
	GetByUserID(ctx context.Context, userID string) ([]models.Notification, error)
	GetByEmail(ctx context.Context, email string) ([]models.Notification, error)
}
