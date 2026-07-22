package kafka

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/anomalyco/hookah-store/notification-service/internal/config"
	"github.com/segmentio/kafka-go"
)

const fc = "notification-service.kafka.consumer"

type EventHandler func(ctx context.Context, payload []byte) error

type Consumer struct {
	cfg      config.KafkaConfig
	consumer *kafka.Reader
	handlers map[string]EventHandler
}

func New(cfg config.KafkaConfig) *Consumer {
	return &Consumer{
		cfg:      cfg,
		handlers: make(map[string]EventHandler),
	}
}

func (c *Consumer) RegisterHandler(topic string, handler EventHandler) {
	c.handlers[topic] = handler
}

func (c *Consumer) Start(ctx context.Context) error {
	if len(c.handlers) == 0 {
		return fmt.Errorf("no handlers registered")
	}

	topics := make([]string, 0, len(c.handlers))
	for topic := range c.handlers {
		topics = append(topics, topic)
	}

	c.consumer = kafka.NewReader(kafka.ReaderConfig{
		Brokers:     c.cfg.Brokers,
		GroupID:     c.cfg.GroupID,
		GroupTopics: topics,
	})

	return c.handleMessages(ctx)
}

func (c *Consumer) handleMessages(ctx context.Context) error {
	for {
		msg, err := c.consumer.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("read message failed: %w", err)
		}

		handler, ok := c.handlers[msg.Topic]
		if !ok {
			slog.Warn("no handler for topic", slog.String("fc", fc), slog.String("topic", msg.Topic))

			continue
		}

		if err := handler(ctx, msg.Value); err != nil {
			slog.Error("handler failed", slog.String("fc", fc), slog.String("topic", msg.Topic), slog.Any("error", err))

			continue
		}

		if err := c.consumer.CommitMessages(ctx, msg); err != nil {
			slog.Error("commit failed", slog.String("fc", fc), slog.String("topic", msg.Topic), slog.Any("error", err))

			continue
		}
	}
}

func (c *Consumer) Close() error {
	if c.consumer == nil {
		return nil
	}

	return c.consumer.Close()
}
