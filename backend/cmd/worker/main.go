package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	kafkago "github.com/segmentio/kafka-go"
	"support-go/backend/internal/notification"
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

	notificationService := notification.NewService(notification.NewLogSender(logger))
	dlqWriter := newWriter(brokers)
	defer dlqWriter.Close()
	retryPolicy := retryPolicy{
		maxAttempts: cfg.NotificationRetryMax,
		backoff:     time.Duration(cfg.NotificationRetryBackoff) * time.Millisecond,
	}

	ticketReader := newReader(brokers, cfg.NotificationConsumerGroup, ticket.TopicTicketEvents)
	commentReader := newReader(brokers, cfg.NotificationConsumerGroup, ticket.TopicCommentEvents)
	defer ticketReader.Close()
	defer commentReader.Close()

	var wg sync.WaitGroup
	errCh := make(chan error, 2)

	wg.Add(2)
	go func() {
		defer wg.Done()
		consumeLoop(ctx, logger, ticketReader, notificationService, dlqWriter, cfg.NotificationDLQTopic, retryPolicy, errCh)
	}()
	go func() {
		defer wg.Done()
		consumeLoop(ctx, logger, commentReader, notificationService, dlqWriter, cfg.NotificationDLQTopic, retryPolicy, errCh)
	}()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		if err != nil && !errors.Is(err, context.Canceled) {
			logger.Error("worker consumer crashed", "error", err)
		}
	case sig := <-shutdownSignal:
		logger.Info("shutdown signal received", "signal", sig.String())
	}

	cancel()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Info("all consumers stopped gracefully")
	case <-time.After(10 * time.Second):
		logger.Warn("consumer shutdown timed out after 10s")
	}

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

func newWriter(brokers []string) *kafkago.Writer {
	return &kafkago.Writer{
		Addr:         kafkago.TCP(brokers...),
		RequiredAcks: kafkago.RequireOne,
		Async:        false,
		Balancer:     &kafkago.LeastBytes{},
	}
}

type retryPolicy struct {
	maxAttempts int
	backoff     time.Duration
}

type deadLetterMessage struct {
	OriginalTopic string              `json:"original_topic"`
	Partition     int                 `json:"partition"`
	Offset        int64               `json:"offset"`
	Attempts      int                 `json:"attempts"`
	FailedAt      time.Time           `json:"failed_at"`
	Error         string              `json:"error"`
	Event         *ticket.DomainEvent `json:"event,omitempty"`
	RawPayload    string              `json:"raw_payload,omitempty"`
}

func consumeLoop(
	ctx context.Context,
	logger *slog.Logger,
	reader *kafkago.Reader,
	notificationService *notification.Service,
	dlqWriter *kafkago.Writer,
	dlqTopic string,
	policy retryPolicy,
	errCh chan<- error,
) {
	for {
		message, err := reader.FetchMessage(ctx)
		if err != nil {
			errCh <- err
			return
		}

		var event ticket.DomainEvent
		if unmarshalErr := json.Unmarshal(message.Value, &event); unmarshalErr != nil {
			logger.Error("failed to decode event", "topic", message.Topic, "partition", message.Partition, "offset", message.Offset, "error", unmarshalErr)
			if dlqErr := writeDeadLetter(ctx, dlqWriter, dlqTopic, message, nil, unmarshalErr, 0); dlqErr != nil {
				errCh <- dlqErr
				return
			}
		} else {
			processErr := processWithRetry(ctx, notificationService, event, policy)
			if processErr != nil {
				logger.Error("notification processing failed, sending to dlq",
					"topic", message.Topic,
					"event_type", event.EventType,
					"entity_id", event.EntityID,
					"event_id", event.ID,
					"error", processErr,
				)
				if dlqErr := writeDeadLetter(ctx, dlqWriter, dlqTopic, message, &event, processErr, policy.maxAttempts); dlqErr != nil {
					errCh <- dlqErr
					return
				}
			}
		}

		if commitErr := reader.CommitMessages(ctx, message); commitErr != nil {
			errCh <- commitErr
			return
		}
	}
}

func processWithRetry(ctx context.Context, service *notification.Service, event ticket.DomainEvent, policy retryPolicy) error {
	maxAttempts := policy.maxAttempts
	if maxAttempts <= 0 {
		maxAttempts = 1
	}

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := service.HandleEvent(ctx, event); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if attempt < maxAttempts && policy.backoff > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(policy.backoff):
			}
		}
	}

	return fmt.Errorf("retries exhausted: %w", lastErr)
}

func writeDeadLetter(
	ctx context.Context,
	writer *kafkago.Writer,
	dlqTopic string,
	message kafkago.Message,
	event *ticket.DomainEvent,
	processErr error,
	attempts int,
) error {
	payload := deadLetterMessage{
		OriginalTopic: message.Topic,
		Partition:     message.Partition,
		Offset:        message.Offset,
		Attempts:      attempts,
		FailedAt:      time.Now().UTC(),
		Error:         processErr.Error(),
		Event:         event,
		RawPayload:    string(message.Value),
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return writer.WriteMessages(ctx, kafkago.Message{
		Topic: dlqTopic,
		Key:   message.Key,
		Value: raw,
		Time:  time.Now().UTC(),
	})
}
