package postgres

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"support-go/backend/internal/ticket"
)

type AuditRepository struct {
	pool *pgxpool.Pool
}

func NewAuditRepository(pool *pgxpool.Pool) *AuditRepository {
	return &AuditRepository{pool: pool}
}

func (repository *AuditRepository) Create(event ticket.TicketEvent) error {
	const query = `
		INSERT INTO ticket_events (
			id, ticket_id, actor_id, event_type, old_value, new_value, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	oldValue, oldErr := marshalNullable(event.OldValue)
	if oldErr != nil {
		return oldErr
	}
	newValue, newErr := marshalNullable(event.NewValue)
	if newErr != nil {
		return newErr
	}

	_, err := repository.pool.Exec(context.Background(), query,
		event.ID,
		event.TicketID,
		event.ActorID,
		event.EventType,
		oldValue,
		newValue,
		event.CreatedAt,
	)

	return err
}

func (repository *AuditRepository) ListByTicketID(ticketID string) ([]ticket.TicketEvent, error) {
	const query = `
		SELECT id, ticket_id, actor_id, event_type, old_value, new_value, created_at
		FROM ticket_events
		WHERE ticket_id = $1
		ORDER BY created_at DESC
	`

	rows, err := repository.pool.Query(context.Background(), query, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]ticket.TicketEvent, 0)
	for rows.Next() {
		var item ticket.TicketEvent
		var oldRaw []byte
		var newRaw []byte
		if scanErr := rows.Scan(
			&item.ID,
			&item.TicketID,
			&item.ActorID,
			&item.EventType,
			&oldRaw,
			&newRaw,
			&item.CreatedAt,
		); scanErr != nil {
			return nil, scanErr
		}

		oldValue, oldErr := unmarshalNullable(oldRaw)
		if oldErr != nil {
			return nil, oldErr
		}
		newValue, newErr := unmarshalNullable(newRaw)
		if newErr != nil {
			return nil, newErr
		}

		item.OldValue = oldValue
		item.NewValue = newValue
		items = append(items, item)
	}

	return items, rows.Err()
}

func marshalNullable(value map[string]any) ([]byte, error) {
	if value == nil {
		return nil, nil
	}

	return json.Marshal(value)
}

func unmarshalNullable(raw []byte) (map[string]any, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	var value map[string]any
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil, err
	}

	return value, nil
}
