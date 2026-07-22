package consumer

import (
	"context"
	"encoding/json"

	"github.com/anomalyco/hookah-store/notification-service/internal/models"
	emailservice "github.com/anomalyco/hookah-store/notification-service/internal/services/email"
	"github.com/anomalyco/hookah-store/notification-service/pkg/kafka"
)

type EmailHandler struct {
	emailService *emailservice.Service
	topic        string
}

func New(service *emailservice.Service, topic string) *EmailHandler {
	return &EmailHandler{
		emailService: service,
		topic:        topic,
	}
}

func (e *EmailHandler) Register(handler kafka.Register) {
	handler.RegisterHandler(e.topic, e.Handle)
}

func (e *EmailHandler) Handle(ctx context.Context, payload []byte) error {
	var event models.Event

	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	return e.emailService.CreateMsg(ctx, &event, event.Type)
}
