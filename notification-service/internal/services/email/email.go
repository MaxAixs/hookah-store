package email

import (
	"context"
	"log/slog"

	"github.com/anomalyco/hookah-store/notification-service/internal/errs"
	"github.com/anomalyco/hookah-store/notification-service/internal/messages"
	"github.com/anomalyco/hookah-store/notification-service/internal/models"
	"github.com/anomalyco/hookah-store/notification-service/internal/repository"
	"github.com/anomalyco/hookah-store/notification-service/pkg/mailgun"
)

type Service struct {
	repo    repository.NotificationRepository
	mailgun mailgun.Mailer
}

func New(repo repository.NotificationRepository, mailgun mailgun.Mailer) *Service {
	return &Service{
		repo:    repo,
		mailgun: mailgun,
	}
}

func (s *Service) CreateMsg(ctx context.Context, event *models.Event) error {
	const fc = "notification_service_service_CreateMsg"

	msg := models.Message{
		To:   event.Payload.Email,
		Name: event.Type,
		Body: messages.MapMsg[event.Type],
	}

	if err := s.mailgun.SendMsg(ctx, msg); err != nil {
		slog.Error("failed to send email", slog.String("fc", fc), slog.Any("error", err))

		return errs.MapErr(err)
	}
	
}
