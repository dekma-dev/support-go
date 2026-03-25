package ticket

import (
	"context"
	"encoding/json"
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

type CommentRepository interface {
	Create(comment Comment) error
	ListByTicketID(ticketID string) ([]Comment, error)
}

type AuditRepository interface {
	Create(event TicketEvent) error
	ListByTicketID(ticketID string) ([]TicketEvent, error)
}

type Service struct {
	repo        Repository
	commentRepo CommentRepository
	auditRepo   AuditRepository
	publisher   EventPublisher
	nowFunc     func() time.Time
	seq         atomic.Uint64
	commentSeq  atomic.Uint64
	eventSeq    atomic.Uint64
	domainSeq   atomic.Uint64
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:    repo,
		nowFunc: time.Now,
	}
}

func NewServiceWithDependencies(repo Repository, commentRepo CommentRepository, auditRepo AuditRepository) *Service {
	return &Service{
		repo:        repo,
		commentRepo: commentRepo,
		auditRepo:   auditRepo,
		nowFunc:     time.Now,
	}
}

func NewServiceWithDependenciesAndPublisher(repo Repository, commentRepo CommentRepository, auditRepo AuditRepository, publisher EventPublisher) *Service {
	return &Service{
		repo:        repo,
		commentRepo: commentRepo,
		auditRepo:   auditRepo,
		publisher:   publisher,
		nowFunc:     time.Now,
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

	_ = service.recordEvent(
		ticket.ID,
		ticket.RequesterID,
		"ticket.created",
		nil,
		map[string]any{
			"status":   ticket.Status,
			"priority": ticket.Priority,
		},
	)
	_ = service.publishEvent(TopicTicketEvents, "ticket.created", ticket.ID, map[string]any{
		"status":       ticket.Status,
		"priority":     ticket.Priority,
		"requester_id": ticket.RequesterID,
	})

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
	before := ticket

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

	_ = service.recordEvent(
		ticket.ID,
		ticket.RequesterID,
		"ticket.updated",
		map[string]any{
			"title":       before.Title,
			"description": before.Description,
			"priority":    before.Priority,
		},
		map[string]any{
			"title":       ticket.Title,
			"description": ticket.Description,
			"priority":    ticket.Priority,
		},
	)
	_ = service.publishEvent(TopicTicketEvents, "ticket.updated", ticket.ID, map[string]any{
		"title":        ticket.Title,
		"description":  ticket.Description,
		"priority":     ticket.Priority,
		"requester_id": ticket.RequesterID,
	})

	return ticket, nil
}

func (service *Service) Assign(id string, assigneeID string, actorID string) (Ticket, error) {
	ticket, err := service.GetByID(id)
	if err != nil {
		return Ticket{}, err
	}
	beforeAssignee := ticket.AssigneeID

	trimmedAssignee := strings.TrimSpace(assigneeID)
	if trimmedAssignee == "" {
		return Ticket{}, fmt.Errorf("%w: assignee_id is required", ErrValidation)
	}

	ticket.AssigneeID = trimmedAssignee
	ticket.UpdatedAt = service.nowFunc().UTC()
	if err := service.repo.Update(ticket); err != nil {
		return Ticket{}, err
	}

	_ = service.recordEvent(
		ticket.ID,
		actorID,
		"ticket.assigned",
		map[string]any{"assignee_id": beforeAssignee},
		map[string]any{"assignee_id": ticket.AssigneeID},
	)
	_ = service.publishEvent(TopicTicketEvents, "ticket.assigned", ticket.ID, map[string]any{
		"assignee_id":  ticket.AssigneeID,
		"requester_id": ticket.RequesterID,
	})

	return ticket, nil
}

func (service *Service) ChangeStatus(id string, status Status, actorID string) (Ticket, error) {
	if err := ValidateStatus(status); err != nil {
		return Ticket{}, err
	}

	ticket, err := service.GetByID(id)
	if err != nil {
		return Ticket{}, err
	}
	beforeStatus := ticket.Status

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

	_ = service.recordEvent(
		ticket.ID,
		actorID,
		"ticket.status.changed",
		map[string]any{"status": beforeStatus},
		map[string]any{"status": ticket.Status},
	)
	_ = service.publishEvent(TopicTicketEvents, "ticket.status.changed", ticket.ID, map[string]any{
		"status":       ticket.Status,
		"requester_id": ticket.RequesterID,
	})

	return ticket, nil
}

type AddCommentInput struct {
	TicketID   string
	AuthorID   string
	Body       string
	IsInternal bool
}

func (service *Service) AddComment(input AddCommentInput) (Comment, error) {
	if service.commentRepo == nil {
		return Comment{}, fmt.Errorf("comment repository is not configured")
	}
	if _, err := service.GetByID(input.TicketID); err != nil {
		return Comment{}, err
	}

	ticketID := strings.TrimSpace(input.TicketID)
	authorID := strings.TrimSpace(input.AuthorID)
	body := strings.TrimSpace(input.Body)
	if ticketID == "" {
		return Comment{}, fmt.Errorf("%w: ticket_id is required", ErrValidation)
	}
	if authorID == "" {
		return Comment{}, fmt.Errorf("%w: author_id is required", ErrValidation)
	}
	if body == "" {
		return Comment{}, fmt.Errorf("%w: body is required", ErrValidation)
	}

	now := service.nowFunc().UTC()
	seq := service.commentSeq.Add(1)
	comment := Comment{
		ID:         fmt.Sprintf("comment-%d", seq),
		TicketID:   ticketID,
		AuthorID:   authorID,
		Body:       body,
		IsInternal: input.IsInternal,
		CreatedAt:  now,
	}

	if err := service.commentRepo.Create(comment); err != nil {
		return Comment{}, err
	}

	_ = service.recordEvent(
		ticketID,
		authorID,
		"comment.added",
		nil,
		map[string]any{"comment_id": comment.ID, "is_internal": comment.IsInternal},
	)
	_ = service.publishEvent(TopicCommentEvents, "comment.added", ticketID, map[string]any{
		"comment_id":  comment.ID,
		"is_internal": comment.IsInternal,
		"author_id":   comment.AuthorID,
	})

	return comment, nil
}

func (service *Service) ListComments(ticketID string, viewerRole Role) ([]Comment, error) {
	if service.commentRepo == nil {
		return nil, fmt.Errorf("comment repository is not configured")
	}
	if _, err := service.GetByID(ticketID); err != nil {
		return nil, err
	}

	comments, err := service.commentRepo.ListByTicketID(strings.TrimSpace(ticketID))
	if err != nil {
		return nil, err
	}

	// Internal comments are visible only to agents/admins.
	if canManageAssignmentsAndStatus(viewerRole) {
		return comments, nil
	}

	filtered := make([]Comment, 0, len(comments))
	for _, item := range comments {
		if !item.IsInternal {
			filtered = append(filtered, item)
		}
	}

	return filtered, nil
}

func (service *Service) ListEvents(ticketID string) ([]TicketEvent, error) {
	if service.auditRepo == nil {
		return nil, fmt.Errorf("audit repository is not configured")
	}
	if _, err := service.GetByID(ticketID); err != nil {
		return nil, err
	}

	return service.auditRepo.ListByTicketID(strings.TrimSpace(ticketID))
}

func (service *Service) recordEvent(ticketID string, actorID string, eventType string, oldValue map[string]any, newValue map[string]any) error {
	if service.auditRepo == nil {
		return nil
	}

	normalizedOld := normalizeJSONMap(oldValue)
	normalizedNew := normalizeJSONMap(newValue)

	event := TicketEvent{
		ID:        fmt.Sprintf("event-%d", service.eventSeq.Add(1)),
		TicketID:  strings.TrimSpace(ticketID),
		ActorID:   strings.TrimSpace(actorID),
		EventType: eventType,
		OldValue:  normalizedOld,
		NewValue:  normalizedNew,
		CreatedAt: service.nowFunc().UTC(),
	}

	return service.auditRepo.Create(event)
}

func normalizeJSONMap(value map[string]any) map[string]any {
	if value == nil {
		return nil
	}

	raw, err := json.Marshal(value)
	if err != nil {
		return value
	}

	var normalized map[string]any
	if unmarshalErr := json.Unmarshal(raw, &normalized); unmarshalErr != nil {
		return value
	}

	return normalized
}

func (service *Service) publishEvent(topic string, eventType string, entityID string, payload map[string]any) error {
	if service.publisher == nil {
		return nil
	}

	event := DomainEvent{
		ID:         fmt.Sprintf("domain-event-%d", service.domainSeq.Add(1)),
		Topic:      topic,
		EventType:  eventType,
		OccurredAt: service.nowFunc().UTC(),
		Producer:   "support-go-api",
		EntityID:   strings.TrimSpace(entityID),
		Payload:    normalizeJSONMap(payload),
	}

	return service.publisher.Publish(context.Background(), event)
}
