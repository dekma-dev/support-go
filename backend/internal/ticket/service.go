package ticket

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"
)

type Repository interface {
	Create(ticket Ticket) error
	GetByID(id string) (Ticket, error)
	List() []Ticket
	Update(ticket Ticket) error
}

type Service struct {
	repo    Repository
	nowFunc func() time.Time
	seq     atomic.Uint64
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:    repo,
		nowFunc: time.Now,
	}
}

type CreateInput struct {
	Title       string
	Description string
	Priority    Priority
	RequesterID string
	SLADueAt    *time.Time
}

type UpdateInput struct {
	Title       *string
	Description *string
	Priority    *Priority
}

func (service *Service) Create(input CreateInput) (Ticket, error) {
	if err := ValidateCreateInput(input.Title, input.Description, input.RequesterID, input.Priority); err != nil {
		return Ticket{}, err
	}

	now := service.nowFunc().UTC()
	sequence := service.seq.Add(1)
	id := fmt.Sprintf("ticket-%d", sequence)

	ticket := Ticket{
		ID:          id,
		PublicID:    fmt.Sprintf("TCK-%06d", sequence),
		Title:       strings.TrimSpace(input.Title),
		Description: strings.TrimSpace(input.Description),
		Status:      StatusNew,
		Priority:    input.Priority,
		RequesterID: strings.TrimSpace(input.RequesterID),
		SLADueAt:    input.SLADueAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := service.repo.Create(ticket); err != nil {
		return Ticket{}, err
	}

	return ticket, nil
}

func (service *Service) List() []Ticket {
	return service.repo.List()
}

func (service *Service) GetByID(id string) (Ticket, error) {
	if strings.TrimSpace(id) == "" {
		return Ticket{}, fmt.Errorf("%w: id is required", ErrValidation)
	}

	return service.repo.GetByID(strings.TrimSpace(id))
}

func (service *Service) Update(id string, input UpdateInput) (Ticket, error) {
	ticket, err := service.GetByID(id)
	if err != nil {
		return Ticket{}, err
	}

	if input.Title != nil {
		trimmed := strings.TrimSpace(*input.Title)
		if trimmed == "" {
			return Ticket{}, fmt.Errorf("%w: title cannot be empty", ErrValidation)
		}
		ticket.Title = trimmed
	}

	if input.Description != nil {
		trimmed := strings.TrimSpace(*input.Description)
		if trimmed == "" {
			return Ticket{}, fmt.Errorf("%w: description cannot be empty", ErrValidation)
		}
		ticket.Description = trimmed
	}

	if input.Priority != nil {
		if err := ValidatePriority(*input.Priority); err != nil {
			return Ticket{}, err
		}
		ticket.Priority = *input.Priority
	}

	ticket.UpdatedAt = service.nowFunc().UTC()
	if err := service.repo.Update(ticket); err != nil {
		return Ticket{}, err
	}

	return ticket, nil
}

func (service *Service) Assign(id string, assigneeID string) (Ticket, error) {
	ticket, err := service.GetByID(id)
	if err != nil {
		return Ticket{}, err
	}

	trimmedAssignee := strings.TrimSpace(assigneeID)
	if trimmedAssignee == "" {
		return Ticket{}, fmt.Errorf("%w: assignee_id is required", ErrValidation)
	}

	ticket.AssigneeID = trimmedAssignee
	ticket.UpdatedAt = service.nowFunc().UTC()
	if err := service.repo.Update(ticket); err != nil {
		return Ticket{}, err
	}

	return ticket, nil
}

func (service *Service) ChangeStatus(id string, status Status) (Ticket, error) {
	if err := ValidateStatus(status); err != nil {
		return Ticket{}, err
	}

	ticket, err := service.GetByID(id)
	if err != nil {
		return Ticket{}, err
	}

	now := service.nowFunc().UTC()
	ticket.Status = status
	ticket.UpdatedAt = now
	if status == StatusClosed {
		ticket.ClosedAt = &now
	} else {
		ticket.ClosedAt = nil
	}

	if err := service.repo.Update(ticket); err != nil {
		return Ticket{}, err
	}

	return ticket, nil
}
