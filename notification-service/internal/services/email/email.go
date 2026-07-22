package email

import (
	"context"
	"fmt"
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

func (s *Service) CreateMsg(ctx context.Context, event *models.Event, eventType string) error {
	const fc = "notification-service.service.CreateMsg"

	if event.Type != models.SignUpEventType && event.Type != models.ResetPasswordEventType {
		slog.Error("unsupported event type %s", event.Type)

		return fmt.Errorf("unknown event type: %s", event.Type)
	}

	notification := &models.Notification{
		ID:        uuid.New(),
		UserID:    event.UserID,
		EventID:   uuid.New(),
		Email:     event.Email,
		EventType: eventType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := s.repo.Create(ctx, notification)
	if err != nil {
		slog.Error("failed to create notification", slog.String("fc", fc), slog.Any("error", err))

		return err
	}

	msg := models.Message{
		To:      event.Email,
		Subject: eventType,
		Body:    messages.MapMsg[eventType],
	}

	msgID, err := s.mailgun.SendMsg(ctx, msg)
	if err != nil {
		slog.Error("failed to send email", slog.String("fc", fc), slog.Any("error", err))

		return err
	}

	if err := s.repo.UpdateMessageID(ctx, id, msgID); err != nil {
		slog.Error("failed to update message ID", slog.String("fc", fc), slog.Any("error", err))

		return err
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
