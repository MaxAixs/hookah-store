package mailgun

import (
	"context"

	"github.com/anomalyco/hookah-store/notification-service/internal/models"
)

type Mailer interface {
	SendMsg(ctx context.Context, msg models.Message) error
}
