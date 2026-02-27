package ticket

import (
	"context"
	"time"
)

const (
	TopicTicketEvents  = "support.ticket.events"
	TopicCommentEvents = "support.comment.events"
)

type DomainEvent struct {
	ID         string         `json:"event_id"`
	Topic      string         `json:"topic"`
	EventType  string         `json:"event_type"`
	OccurredAt time.Time      `json:"occurred_at"`
	Producer   string         `json:"producer"`
	EntityID   string         `json:"entity_id"`
	TraceID    string         `json:"trace_id,omitempty"`
	Payload    map[string]any `json:"payload,omitempty"`
}

type EventPublisher interface {
	Publish(ctx context.Context, event DomainEvent) error
}
