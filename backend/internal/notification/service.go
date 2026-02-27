package notification

import (
	"context"
	"fmt"
	"log/slog"

	"support-go/backend/internal/ticket"
)

type Message struct {
	Recipient string
	Subject   string
	Body      string
}

type Sender interface {
	Send(ctx context.Context, message Message) error
}

type LogSender struct {
	logger *slog.Logger
}

func NewLogSender(logger *slog.Logger) *LogSender {
	return &LogSender{logger: logger}
}

func (sender *LogSender) Send(_ context.Context, message Message) error {
	sender.logger.Info("notification sent",
		"recipient", message.Recipient,
		"subject", message.Subject,
		"body", message.Body,
	)
	return nil
}

type Service struct {
	sender Sender
}

func NewService(sender Sender) *Service {
	return &Service{sender: sender}
}

func (service *Service) HandleEvent(ctx context.Context, event ticket.DomainEvent) error {
	message, ok := mapEventToMessage(event)
	if !ok {
		return nil
	}

	return service.sender.Send(ctx, message)
}

func mapEventToMessage(event ticket.DomainEvent) (Message, bool) {
	switch event.EventType {
	case "ticket.created":
		return Message{
			Recipient: valueOrFallback(event.Payload, "requester_id", "requester"),
			Subject:   "Ticket created",
			Body:      fmt.Sprintf("Ticket %s was created", event.EntityID),
		}, true
	case "ticket.updated":
		return Message{
			Recipient: valueOrFallback(event.Payload, "requester_id", "requester"),
			Subject:   "Ticket updated",
			Body:      fmt.Sprintf("Ticket %s was updated", event.EntityID),
		}, true
	case "ticket.assigned":
		return Message{
			Recipient: valueOrFallback(event.Payload, "assignee_id", "agent"),
			Subject:   "Ticket assigned",
			Body:      fmt.Sprintf("Ticket %s was assigned to you", event.EntityID),
		}, true
	case "ticket.status.changed":
		status := valueOrFallback(event.Payload, "status", "updated")
		return Message{
			Recipient: valueOrFallback(event.Payload, "requester_id", "requester"),
			Subject:   "Ticket status changed",
			Body:      fmt.Sprintf("Ticket %s status changed to %s", event.EntityID, status),
		}, true
	case "comment.added":
		return Message{
			Recipient: valueOrFallback(event.Payload, "author_id", "requester"),
			Subject:   "New comment",
			Body:      fmt.Sprintf("New comment added for ticket %s", event.EntityID),
		}, true
	default:
		return Message{}, false
	}
}

func valueOrFallback(payload map[string]any, key string, fallback string) string {
	if payload == nil {
		return fallback
	}

	raw, ok := payload[key]
	if !ok {
		return fallback
	}

	value, ok := raw.(string)
	if !ok || value == "" {
		return fallback
	}

	return value
}
