package memory

import (
	"sort"
	"sync"

	"support-go/backend/internal/ticket"
)

type Repository struct {
	mu      sync.RWMutex
	tickets map[string]ticket.Ticket
}

func NewRepository() *Repository {
	return &Repository{
		tickets: make(map[string]ticket.Ticket),
	}
}

func (repository *Repository) Create(value ticket.Ticket) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	repository.tickets[value.ID] = value
	return nil
}

func (repository *Repository) GetByID(id string) (ticket.Ticket, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	value, ok := repository.tickets[id]
	if !ok {
		return ticket.Ticket{}, ticket.ErrNotFound
	}

	return value, nil
}

func (repository *Repository) List() []ticket.Ticket {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	items := make([]ticket.Ticket, 0, len(repository.tickets))
	for _, value := range repository.tickets {
		items = append(items, value)
	}

	sort.Slice(items, func(left, right int) bool {
		return items[left].CreatedAt.After(items[right].CreatedAt)
	})

	return items
}

func (repository *Repository) Update(value ticket.Ticket) error {
	repository.mu.Lock()
	defer repository.mu.Unlock()

	if _, ok := repository.tickets[value.ID]; !ok {
		return ticket.ErrNotFound
	}

	repository.tickets[value.ID] = value
	return nil
}
