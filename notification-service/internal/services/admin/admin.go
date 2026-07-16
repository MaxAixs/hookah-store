package admin

import (
	"context"
	"log/slog"

	"github.com/anomalyco/hookah-store/notification-service/internal/errs"
	"github.com/anomalyco/hookah-store/notification-service/internal/models"
	"github.com/anomalyco/hookah-store/notification-service/internal/repository"
)

type Service struct {
	notifRepo repository.NotificationRepository
}

func New(notifRepo repository.NotificationRepository) *Service {
	return &Service{notifRepo: notifRepo}
}

func (s *Service) GetByUserID(ctx context.Context, userID string) ([]models.NotificationResponse, error) {
	const fc = "notification-service.services.GetByUserID"

	notifications, err := s.notifRepo.GetByUserID(ctx, userID)
	if err != nil {
		slog.Error("failed to get notifications by user id", slog.String("fc", fc), slog.Any("error", err))

		return nil, errs.MapErr(err)
	}

	return toResponses(notifications), nil
}

func (s *Service) GetByEmail(ctx context.Context, email string) ([]models.NotificationResponse, error) {
	const fc = "notification-service.services.GetByEmail"

	notifications, err := s.notifRepo.GetByEmail(ctx, email)
	if err != nil {
		slog.Error("failed to get notifications by email", slog.String("fc", fc), slog.Any("error", err))

		return nil, errs.MapErr(err)
	}

	return toResponses(notifications), nil
}

func toResponses(notifications []models.Notification) []models.NotificationResponse {
	resp := make([]models.NotificationResponse, 0, len(notifications))
	for _, n := range notifications {
		resp = append(resp, models.NotificationResponse{
			ID:        n.ID,
			UserID:    n.UserID,
			Email:     n.Email,
			EventType: n.EventType,
			Status:    n.Status,
			CreatedAt: n.CreatedAt,
			UpdatedAt: n.UpdatedAt,
		})
	}
	return resp
}
