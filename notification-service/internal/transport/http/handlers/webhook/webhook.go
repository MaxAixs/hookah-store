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
	var mailgunEvent models.MailgunEvent

	if err := ctx.ShouldBind(&mailgunEvent); err != nil {
		ctx.Status(http.StatusBadRequest)

		return
	}

	if err := h.EmailService.UpdateStatus(ctx, mailgunEvent); err != nil {
		ctx.Status(http.StatusInternalServerError)

		return
	}

	ctx.Status(http.StatusOK)
}
