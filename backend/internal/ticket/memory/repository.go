package memory

import (
	"context"
	"sort"
	"strings"
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

func (repository *Repository) ListWithFilter(_ context.Context, filter ticket.ListFilter) ([]ticket.Ticket, int, error) {
	repository.mu.RLock()
	defer repository.mu.RUnlock()

	matched := make([]ticket.Ticket, 0, len(repository.tickets))
	for _, value := range repository.tickets {
		if !matchesFilter(value, filter) {
			continue
		}
		matched = append(matched, value)
	}

	switch filter.Sort {
	case ticket.SortCreatedAtAsc:
		sort.Slice(matched, func(l, r int) bool { return matched[l].CreatedAt.Before(matched[r].CreatedAt) })
	case ticket.SortUpdatedAtDesc:
		sort.Slice(matched, func(l, r int) bool { return matched[l].UpdatedAt.After(matched[r].UpdatedAt) })
	case ticket.SortUpdatedAtAsc:
		sort.Slice(matched, func(l, r int) bool { return matched[l].UpdatedAt.Before(matched[r].UpdatedAt) })
	default:
		sort.Slice(matched, func(l, r int) bool { return matched[l].CreatedAt.After(matched[r].CreatedAt) })
	}

	total := len(matched)
	if filter.Offset >= total {
		return []ticket.Ticket{}, total, nil
	}

	end := filter.Offset + filter.Limit
	if end > total {
		end = total
	}

	return matched[filter.Offset:end], total, nil
}

func matchesFilter(value ticket.Ticket, filter ticket.ListFilter) bool {
	if len(filter.Statuses) > 0 {
		found := false
		for _, s := range filter.Statuses {
			if value.Status == s {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if len(filter.Priorities) > 0 {
		found := false
		for _, p := range filter.Priorities {
			if value.Priority == p {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if filter.AssigneeID != "" && value.AssigneeID != filter.AssigneeID {
		return false
	}

	if filter.Search != "" {
		q := strings.ToLower(filter.Search)
		if !strings.Contains(strings.ToLower(value.Title), q) &&
			!strings.Contains(strings.ToLower(value.Description), q) {
			return false
		}
	}

	return true
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
