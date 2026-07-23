package email

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
	repo       repository.NotificationRepository
	signingKey string
	mailgun    mailgun.Mailer
}

func New(repo repository.NotificationRepository, signingKey string, mailgun mailgun.Mailer) *Service {
	return &Service{
		repo:       repo,
		signingKey: signingKey,
		mailgun:    mailgun,
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

	slog.Info("notification created", slog.String("msgID", msgID))

	if err := s.repo.UpdateMessageID(ctx, id, msgID); err != nil {
		slog.Error("failed to update message ID", slog.String("fc", fc), slog.Any("error", err))

		return err
	}

	return nil
}

func (s *Service) UpdateStatus(ctx context.Context, mailgunData models.MailgunWebhook) error {
	const fc = "notification-service.service.UpdateStatus"

	if err := verifySignature(s.signingKey, mailgunData.Signature.Timestamp, mailgunData.Signature.Token,
		mailgunData.Signature.Signature); err != nil {
		slog.Error("failed to verify signature", slog.String("fc", fc), slog.Any("error", err))

		return err
	}


	if err := s.repo.UpdateStatus(ctx, mailgunData.EventData.Message.Headers.To,
		mailgunData.EventData.Message.Headers.MessageID, models.MapMsgStatus[mailgunData.EventData.Event]); err != nil {
		slog.Error("failed to update status", slog.String("fc", fc), slog.Any("error", err))

		return errs.MapErr(err)
	}

	return nil
}

func verifySignature(signingKey, timestamp, token, signature string) error {
	mac := hmac.New(sha256.New, []byte(signingKey))
	mac.Write([]byte(timestamp + token))
	expected := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expected), []byte(signature)) {
		return errs.ErrInvalidSignature
	}

	return nil
}
