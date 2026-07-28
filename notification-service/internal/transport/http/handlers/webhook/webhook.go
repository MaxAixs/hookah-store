package webhook

import (
	"net/http"

	"github.com/anomalyco/hookah-store/notification-service/internal/models"
	"github.com/anomalyco/hookah-store/notification-service/internal/services/email"
	"github.com/gin-gonic/gin"
)

type WebhookHandlers struct {
	EmailService *email.EmailService
}

func New(emailService *email.EmailService) *WebhookHandlers {
	return &WebhookHandlers{
		EmailService: emailService,
	}
}

func (h *WebhookHandlers) Register(router *gin.RouterGroup) {
	router.POST("/webhook/mailgun", h.MailgunWebhook)
}

func (h *WebhookHandlers) ShutDown() {}

func (h *WebhookHandlers) MailgunWebhook(ctx *gin.Context) {
	var mailgunData models.MailgunWebhook
	if err := ctx.ShouldBind(&mailgunData); err != nil {
		ctx.Status(http.StatusBadRequest)
	}

	if err := h.EmailService.UpdateStatus(ctx, mailgunData); err != nil {
		ctx.Status(http.StatusInternalServerError)

		return
	}

	ctx.Status(http.StatusOK)
}
