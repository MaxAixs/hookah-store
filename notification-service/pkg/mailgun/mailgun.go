package mailgun

import (
	"context"
	"fmt"

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
	newMsg := mailgun.NewMessage(c.cfg.Domain, c.cfg.From, msg.Subject, msg.Body, msg.To)
	if newMsg == nil {
		return "", fmt.Errorf("new Message is nil")
	}

	resp, err := c.mg.Send(ctx, newMsg)
	if err != nil {
		return "", fmt.Errorf("failed to send email: %w", err)
	}

	return resp.ID, nil
}
