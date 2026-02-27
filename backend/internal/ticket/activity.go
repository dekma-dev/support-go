package ticket

import "time"

type Comment struct {
	ID         string    `json:"id"`
	TicketID   string    `json:"ticket_id"`
	AuthorID   string    `json:"author_id"`
	Body       string    `json:"body"`
	IsInternal bool      `json:"is_internal"`
	CreatedAt  time.Time `json:"created_at"`
}

type TicketEvent struct {
	ID        string         `json:"id"`
	TicketID  string         `json:"ticket_id"`
	ActorID   string         `json:"actor_id"`
	EventType string         `json:"event_type"`
	OldValue  map[string]any `json:"old_value,omitempty"`
	NewValue  map[string]any `json:"new_value,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}
