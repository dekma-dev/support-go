package notification

import (
	"context"
	"testing"

	"support-go/backend/internal/ticket"
)

type testSender struct {
	messages []Message
}

func (sender *testSender) Send(_ context.Context, message Message) error {
	sender.messages = append(sender.messages, message)
	return nil
}

func TestHandleEventCreatesMessage(t *testing.T) {
	sender := &testSender{}
	service := NewService(sender)

	err := service.HandleEvent(context.Background(), ticket.DomainEvent{
		EventType: "ticket.created",
		EntityID:  "ticket-1",
		Payload: map[string]any{
			"requester_id": "user-1",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(sender.messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(sender.messages))
	}
	if sender.messages[0].Recipient != "user-1" {
		t.Fatalf("expected recipient user-1, got %s", sender.messages[0].Recipient)
	}
}

func TestHandleEventIgnoresUnknownType(t *testing.T) {
	sender := &testSender{}
	service := NewService(sender)

	err := service.HandleEvent(context.Background(), ticket.DomainEvent{
		EventType: "unknown.event",
		EntityID:  "ticket-1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(sender.messages) != 0 {
		t.Fatalf("expected no messages, got %d", len(sender.messages))
	}
}
