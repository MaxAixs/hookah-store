package mailgun

import (
	"context"
	"fmt"
	"strings"

	"github.com/anomalyco/hookah-store/notification-service/internal/config"
	"github.com/anomalyco/hookah-store/notification-service/internal/models"
	"github.com/mailgun/mailgun-go/v5"
)

type Client struct {
	cfg config.MailGunConfig
	mg  *mailgun.Client
}

func New(cfg config.MailGunConfig) *Client {
	return &Client{
		cfg: cfg,
		mg:  mailgun.NewMailgun(cfg.APIKey),
	}
}

func (c *Client) SendMsg(ctx context.Context, msg models.Message) (string, error) {
	switch {
	case c.cfg.Domain == "" || c.cfg.From == "":
		return "", fmt.Errorf("mailgun config is not set: domain=%q from=%q", c.cfg.Domain, c.cfg.From)
	case msg.To == "":
		return "", fmt.Errorf("recipient email is empty")
	case msg.Body == "":
		return "", fmt.Errorf("message body is empty")
	}

	newMsg := mailgun.NewMessage(c.cfg.Domain, c.cfg.From, msg.Subject, msg.Body, msg.To)

	resp, err := c.mg.Send(ctx, newMsg)
	if err != nil {
		return "", fmt.Errorf("failed to send email: %w", err)
	}

	newID := strings.Trim(resp.ID, "<>")

	return newID, nil
}
