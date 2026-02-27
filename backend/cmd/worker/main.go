package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	kafkago "github.com/segmentio/kafka-go"
	"support-go/backend/internal/platform/config"
	platformkafka "support-go/backend/internal/platform/kafka"
	"support-go/backend/internal/ticket"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))

	brokers := platformkafka.ParseBrokers(cfg.KafkaBrokers)
	if len(brokers) == 0 {
		logger.Error("KAFKA_BROKERS is required for worker")
		os.Exit(1)
	}

	logger.Info("starting notification worker", "brokers", brokers, "group", cfg.NotificationConsumerGroup)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticketReader := newReader(brokers, cfg.NotificationConsumerGroup, ticket.TopicTicketEvents)
	commentReader := newReader(brokers, cfg.NotificationConsumerGroup, ticket.TopicCommentEvents)
	defer ticketReader.Close()
	defer commentReader.Close()

	errCh := make(chan error, 2)
	go consumeLoop(ctx, logger, ticketReader, errCh)
	go consumeLoop(ctx, logger, commentReader, errCh)

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, context.Canceled) {
			logger.Error("worker consumer crashed", "error", err)
			os.Exit(1)
		}
	case sig := <-shutdownSignal:
		logger.Info("shutdown signal received", "signal", sig.String())
	}

	cancel()
	time.Sleep(500 * time.Millisecond)
	logger.Info("notification worker stopped")
}

func newReader(brokers []string, groupID string, topic string) *kafkago.Reader {
	return kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		Topic:       topic,
		StartOffset: kafkago.LastOffset,
		MinBytes:    1,
		MaxBytes:    10e6,
		MaxWait:     500 * time.Millisecond,
	})
}

func consumeLoop(ctx context.Context, logger *slog.Logger, reader *kafkago.Reader, errCh chan<- error) {
	for {
		message, err := reader.FetchMessage(ctx)
		if err != nil {
			errCh <- err
			return
		}

		var event ticket.DomainEvent
		if unmarshalErr := json.Unmarshal(message.Value, &event); unmarshalErr != nil {
			logger.Error("failed to decode event", "topic", message.Topic, "partition", message.Partition, "offset", message.Offset, "error", unmarshalErr)
		} else {
			logger.Info("notification event received",
				"topic", message.Topic,
				"event_type", event.EventType,
				"entity_id", event.EntityID,
				"event_id", event.ID,
			)
		}

		if commitErr := reader.CommitMessages(ctx, message); commitErr != nil {
			errCh <- commitErr
			return
		}
	}
}
