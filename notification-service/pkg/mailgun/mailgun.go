package mailgun

import (
	"github.com/anomalyco/hookah-store/notification-service/internal/config"
	"github.com/anomalyco/hookah-store/notification-service/internal/models"
	"github.com/mailgun/mailgun-go/v5"
)

type Client struct {
	mg     *mailgun.Client
	domain string
}

func New(cfg config.MailGunConfig) *Client {
	return &Client{
		mg:     mailgun.NewMailgun(cfg.APIKey),
		domain: cfg.Domain,
	}
}

func (c *Client) SendMsg(msg models.Notification) error {
	c.mg.
}
