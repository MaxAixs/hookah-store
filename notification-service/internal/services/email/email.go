package email

import (
	"context"
	"log/slog"
	"time"

	"github.com/anomalyco/hookah-store/notification-service/internal/errs"
	"github.com/anomalyco/hookah-store/notification-service/internal/messages"
	"github.com/anomalyco/hookah-store/notification-service/internal/models"
	"github.com/anomalyco/hookah-store/notification-service/internal/repository"
	"github.com/anomalyco/hookah-store/notification-service/pkg/mailgun"
	"github.com/google/uuid"
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
	const fc = "notification-service.service.CreateMsg"

	notification := &models.Notification{
		ID:        uuid.New(),
		UserID:    event.Payload.UserID,
		EventID:   event.ID,
		Email:     event.Payload.Email,
		EventType: event.Type,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := s.repo.Create(ctx, notification)
	if err != nil {
		slog.Error("failed to create notification", slog.String("fc", fc), slog.Any("error", err))

		return errs.MapErr(err)
	}

	msg := models.Message{
		To:      event.Payload.Email,
		Subject: event.Type,
		Body:    messages.MapMsg[event.Type],
	}

	msgID, err := s.mailgun.SendMsg(ctx, msg)
	if err != nil {
		slog.Error("failed to send email", slog.String("fc", fc), slog.Any("error", err))

		return errs.MapErr(err)
	}

	if err := s.repo.UpdateMessageID(ctx, id, msgID); err != nil {
		slog.Error("failed to update message ID", slog.String("fc", fc), slog.Any("error", err))

		return errs.MapErr(err)
	}

	return nil
}

func (s *Service) UpdateStatus(ctx context.Context, mailgunEvent models.MailgunEvent) error {
	const fc = "notification-service.service.UpdateStatus"

	if err := s.repo.UpdateStatus(ctx, mailgunEvent.ID, models.MapMsgStatus[mailgunEvent.Event]); err != nil {
		slog.Error("failed to update status", slog.String("fc", fc), slog.Any("error", err))

		return errs.MapErr(err)
	}

	return nil
}
