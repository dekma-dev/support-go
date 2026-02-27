package kafka

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	kafkago "github.com/segmentio/kafka-go"
	"support-go/backend/internal/ticket"
)

type Publisher struct {
	writer *kafkago.Writer
}

type NoopPublisher struct{}

func NewPublisher(brokers []string) *Publisher {
	return &Publisher{
		writer: &kafkago.Writer{
			Addr:         kafkago.TCP(brokers...),
			RequiredAcks: kafkago.RequireOne,
			Async:        true,
			Balancer:     &kafkago.LeastBytes{},
		},
	}
}

func NewNoopPublisher() *NoopPublisher {
	return &NoopPublisher{}
}

func ParseBrokers(raw string) []string {
	parts := strings.Split(raw, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			items = append(items, trimmed)
		}
	}
	return items
}

func (publisher *Publisher) Publish(ctx context.Context, event ticket.DomainEvent) error {
	raw, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return publisher.writer.WriteMessages(ctx, kafkago.Message{
		Topic: event.Topic,
		Key:   []byte(event.EntityID),
		Value: raw,
		Time:  time.Now().UTC(),
	})
}

func (publisher *Publisher) Close() error {
	return publisher.writer.Close()
}

func (publisher *NoopPublisher) Publish(_ context.Context, _ ticket.DomainEvent) error {
	return nil
}
