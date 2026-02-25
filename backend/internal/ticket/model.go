package ticket

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Status string

type Priority string

const (
	StatusNew             Status = "new"
	StatusOpen            Status = "open"
	StatusPendingCustomer Status = "pending_customer"
	StatusPendingInternal Status = "pending_internal"
	StatusResolved        Status = "resolved"
	StatusClosed          Status = "closed"
)

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
	PriorityUrgent Priority = "urgent"
)

var (
	ErrNotFound   = errors.New("ticket not found")
	ErrValidation = errors.New("validation failed")
)

type Ticket struct {
	ID          string     `json:"id"`
	PublicID    string     `json:"public_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	Priority    Priority   `json:"priority"`
	RequesterID string     `json:"requester_id"`
	AssigneeID  string     `json:"assignee_id,omitempty"`
	SLADueAt    *time.Time `json:"sla_due_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ClosedAt    *time.Time `json:"closed_at,omitempty"`
}

func ValidateStatus(value Status) error {
	switch value {
	case StatusNew, StatusOpen, StatusPendingCustomer, StatusPendingInternal, StatusResolved, StatusClosed:
		return nil
	default:
		return fmt.Errorf("%w: unsupported status %q", ErrValidation, value)
	}
}

func ValidatePriority(value Priority) error {
	switch value {
	case PriorityLow, PriorityMedium, PriorityHigh, PriorityUrgent:
		return nil
	default:
		return fmt.Errorf("%w: unsupported priority %q", ErrValidation, value)
	}
}

func ValidateCreateInput(title, description, requesterID string, priority Priority) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("%w: title is required", ErrValidation)
	}
	if strings.TrimSpace(description) == "" {
		return fmt.Errorf("%w: description is required", ErrValidation)
	}
	if strings.TrimSpace(requesterID) == "" {
		return fmt.Errorf("%w: requester_id is required", ErrValidation)
	}

	return ValidatePriority(priority)
}
