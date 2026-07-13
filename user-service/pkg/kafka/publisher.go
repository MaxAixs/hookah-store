package kafka

import (
	"context"
	"fmt"

	"github.com/anomalyco/hookah-store/user-service/internal/config"
	"github.com/segmentio/kafka-go"
)

type Publisher struct {
	writer kafka.Writer
}

func NewPublisher(cfg config.KafkaCfg) *Publisher {
	return &Publisher{
		writer: kafka.Writer{
			Addr:         kafka.TCP(cfg.Brokers...),
			Balancer:     &kafka.Hash{},
			RequiredAcks: kafka.RequiredAcks(cfg.RequiredAcks),
			Async:        cfg.Async,
		},
	}
}

type Message struct {
	Topic string
	Key   string
	Value []byte
}

func (p *Publisher) Publish(ctx context.Context, msg Message) error {
	err := p.writer.WriteMessages(ctx, kafka.Message{
		Topic: msg.Topic,
		Key:   []byte(msg.Key),
		Value: msg.Value,
	})
	if err != nil {
		return fmt.Errorf("kafka publish failed: %w", err)
	}

	return nil
}

func (p *Publisher) Close() error {
	return p.writer.Close()
}
