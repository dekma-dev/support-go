package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"support-go/backend/internal/notification"
	"support-go/backend/internal/ticket"
)

type flakySender struct {
	failuresBeforeSuccess int
	attempts              int
}

func (sender *flakySender) Send(_ context.Context, _ notification.Message) error {
	sender.attempts++
	if sender.attempts <= sender.failuresBeforeSuccess {
		return errors.New("temporary send error")
	}
	return nil
}

func TestProcessWithRetryEventuallySucceeds(t *testing.T) {
	sender := &flakySender{failuresBeforeSuccess: 2}
	service := notification.NewService(sender)
	err := processWithRetry(context.Background(), service, ticket.DomainEvent{
		EventType: "ticket.created",
		EntityID:  "ticket-1",
	}, retryPolicy{maxAttempts: 3, backoff: 1 * time.Millisecond})

	if err != nil {
		t.Fatalf("expected success after retries, got error: %v", err)
	}
	if sender.attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", sender.attempts)
	}
}

func TestProcessWithRetryFailsAfterExhaustedAttempts(t *testing.T) {
	sender := &flakySender{failuresBeforeSuccess: 10}
	service := notification.NewService(sender)
	err := processWithRetry(context.Background(), service, ticket.DomainEvent{
		EventType: "ticket.created",
		EntityID:  "ticket-1",
	}, retryPolicy{maxAttempts: 3, backoff: 1 * time.Millisecond})

	if err == nil {
		t.Fatal("expected retries exhausted error, got nil")
	}
	if sender.attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", sender.attempts)
	}
}
