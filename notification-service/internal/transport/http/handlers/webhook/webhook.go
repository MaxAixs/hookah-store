package webhook

import (
	"net/http"

	"github.com/anomalyco/hookah-store/notification-service/internal/models"
	"github.com/anomalyco/hookah-store/notification-service/internal/services/email"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	EmailService *email.Service
}

func New(emailService *email.Service) *Handlers {
	return &Handlers{
		EmailService: emailService,
	}
}

func (h *Handlers) Register(router *gin.RouterGroup) {
	router.POST("/webhook/mailgun", h.MailgunWebhook)
}

func (h *Handlers) ShutDown() {}

func (h *Handlers) MailgunWebhook(ctx *gin.Context) {
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
