package relay

import (
	"context"
	"log/slog"
	"time"

	"github.com/anomalyco/hookah-store/user-service/internal/repository"
	"github.com/anomalyco/hookah-store/user-service/pkg/kafka"
)

type OutboxRelay struct {
	outBoxRepo repository.OutBoxRepository
	publisher  *kafka.Publisher
}

func NewOutboxRelay(repo repository.OutBoxRepository, publisher *kafka.Publisher) *OutboxRelay {
	return &OutboxRelay{
		outBoxRepo: repo,
		publisher:  publisher,
	}
}

func (r *OutboxRelay) Run(ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := r.publishBatch(ctx); err != nil {
				slog.Error("outbox relay batch failed", slog.Any("error", err))
			}
		}
	}
}

func (r *OutboxRelay) publishBatch(ctx context.Context) error {
	events, err := r.outBoxRepo.FetchUnpublishedEvents(10)
	if err != nil {
		slog.Error("failed to get unpublished events", slog.Any("error", err))

		return err
	}

	for _, e := range events {
		err := r.publisher.Publish(ctx, kafka.Message{
			Topic: e.Topic,
			Key:   e.Key,
			Value: e.Payload,
		})
		if err != nil {
			slog.Error("failed to publish event", slog.Any("error", err))

			continue
		}
		if err := r.outBoxRepo.MarkPublishedEvents(ctx, e.ID); err != nil {
			slog.Error("failed to mark published event", slog.Any("error", err))

			return err
		}
	}

	return nil
}
