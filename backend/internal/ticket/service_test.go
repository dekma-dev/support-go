package ticket_test

import (
	"testing"

	"support-go/backend/internal/ticket"
	"support-go/backend/internal/ticket/memory"
)

func TestServiceCreateAndGet(t *testing.T) {
	repository := memory.NewRepository()
	service := ticket.NewService(repository)

	created, err := service.Create(ticket.CreateInput{
		Title:       "Email delivery failed",
		Description: "Customer does not receive password reset mail",
		Priority:    ticket.PriorityHigh,
		RequesterID: "user-1",
	})
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}

	if created.Status != ticket.StatusNew {
		t.Fatalf("expected status %s, got %s", ticket.StatusNew, created.Status)
	}

	loaded, err := service.GetByID(created.ID)
	if err != nil {
		t.Fatalf("unexpected get error: %v", err)
	}

	if loaded.ID != created.ID {
		t.Fatalf("expected id %s, got %s", created.ID, loaded.ID)
	}
}

func TestServiceAssignAndClose(t *testing.T) {
	repository := memory.NewRepository()
	service := ticket.NewService(repository)

	created, err := service.Create(ticket.CreateInput{
		Title:       "Access denied",
		Description: "Agent cannot open dashboard",
		Priority:    ticket.PriorityMedium,
		RequesterID: "user-2",
	})
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}

	assigned, err := service.Assign(created.ID, "agent-1", "admin-1")
	if err != nil {
		t.Fatalf("unexpected assign error: %v", err)
	}

	if assigned.AssigneeID != "agent-1" {
		t.Fatalf("expected assignee agent-1, got %s", assigned.AssigneeID)
	}

	closed, err := service.ChangeStatus(created.ID, ticket.StatusClosed, "admin-1")
	if err != nil {
		t.Fatalf("unexpected status error: %v", err)
	}

	if closed.ClosedAt == nil {
		t.Fatal("expected closed_at to be set")
	}
}
