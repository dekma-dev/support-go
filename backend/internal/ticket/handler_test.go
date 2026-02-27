package ticket_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"

	"support-go/backend/internal/ticket"
	"support-go/backend/internal/ticket/memory"
)

type testCommentRepo struct {
	items []ticket.Comment
}

func (repository *testCommentRepo) Create(comment ticket.Comment) error {
	repository.items = append(repository.items, comment)
	return nil
}

func (repository *testCommentRepo) ListByTicketID(ticketID string) ([]ticket.Comment, error) {
	result := make([]ticket.Comment, 0)
	for _, item := range repository.items {
		if item.TicketID == ticketID {
			result = append(result, item)
		}
	}
	sort.Slice(result, func(left, right int) bool {
		return result[left].CreatedAt.After(result[right].CreatedAt)
	})
	return result, nil
}

type testAuditRepo struct {
	items []ticket.TicketEvent
}

func (repository *testAuditRepo) Create(event ticket.TicketEvent) error {
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}
	repository.items = append(repository.items, event)
	return nil
}

func (repository *testAuditRepo) ListByTicketID(ticketID string) ([]ticket.TicketEvent, error) {
	result := make([]ticket.TicketEvent, 0)
	for _, item := range repository.items {
		if item.TicketID == ticketID {
			result = append(result, item)
		}
	}
	sort.Slice(result, func(left, right int) bool {
		return result[left].CreatedAt.After(result[right].CreatedAt)
	})
	return result, nil
}

func setupTicketRouter() *http.ServeMux {
	repository := memory.NewRepository()
	commentRepo := &testCommentRepo{}
	auditRepo := &testAuditRepo{}
	service := ticket.NewServiceWithDependencies(repository, commentRepo, auditRepo)
	mux := http.NewServeMux()
	ticket.RegisterRoutes(mux, service)
	return mux
}

func createTicketForHandlerTest(t *testing.T, mux *http.ServeMux) string {
	t.Helper()

	body := map[string]any{
		"title":        "Login issue",
		"description":  "Cannot login to dashboard",
		"priority":     "high",
		"requester_id": "user-1",
	}

	payload, _ := json.Marshal(body)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/tickets", bytes.NewReader(payload))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201 on create, got %d", recorder.Code)
	}

	var created map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &created); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	id, _ := created["id"].(string)
	if id == "" {
		t.Fatal("expected ticket id in response")
	}

	return id
}

func TestAssignRequiresAgentOrAdminRole(t *testing.T) {
	mux := setupTicketRouter()
	id := createTicketForHandlerTest(t, mux)

	payload := []byte(`{"assignee_id":"agent-1"}`)

	requestWithoutRole := httptest.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id+"/assign", bytes.NewReader(payload))
	requestWithoutRole.Header.Set("Content-Type", "application/json")
	recorderWithoutRole := httptest.NewRecorder()
	mux.ServeHTTP(recorderWithoutRole, requestWithoutRole)
	if recorderWithoutRole.Code != http.StatusForbidden {
		t.Fatalf("expected 403 without role, got %d", recorderWithoutRole.Code)
	}

	requestClientRole := httptest.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id+"/assign", bytes.NewReader(payload))
	requestClientRole.Header.Set("Content-Type", "application/json")
	requestClientRole.Header.Set("X-User-Role", "client")
	recorderClientRole := httptest.NewRecorder()
	mux.ServeHTTP(recorderClientRole, requestClientRole)
	if recorderClientRole.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for client role, got %d", recorderClientRole.Code)
	}

	requestAgentRole := httptest.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id+"/assign", bytes.NewReader(payload))
	requestAgentRole.Header.Set("Content-Type", "application/json")
	requestAgentRole.Header.Set("X-User-Role", "agent")
	recorderAgentRole := httptest.NewRecorder()
	mux.ServeHTTP(recorderAgentRole, requestAgentRole)
	if recorderAgentRole.Code != http.StatusOK {
		t.Fatalf("expected 200 for agent role, got %d", recorderAgentRole.Code)
	}
}

func TestStatusRequiresAgentOrAdminRole(t *testing.T) {
	mux := setupTicketRouter()
	id := createTicketForHandlerTest(t, mux)

	payload := []byte(`{"status":"resolved"}`)

	requestClientRole := httptest.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id+"/status", bytes.NewReader(payload))
	requestClientRole.Header.Set("Content-Type", "application/json")
	requestClientRole.Header.Set("X-User-Role", "client")
	recorderClientRole := httptest.NewRecorder()
	mux.ServeHTTP(recorderClientRole, requestClientRole)
	if recorderClientRole.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for client role, got %d", recorderClientRole.Code)
	}

	requestAdminRole := httptest.NewRequest(http.MethodPatch, "/api/v1/tickets/"+id+"/status", bytes.NewReader(payload))
	requestAdminRole.Header.Set("Content-Type", "application/json")
	requestAdminRole.Header.Set("X-User-Role", "admin")
	recorderAdminRole := httptest.NewRecorder()
	mux.ServeHTTP(recorderAdminRole, requestAdminRole)
	if recorderAdminRole.Code != http.StatusOK {
		t.Fatalf("expected 200 for admin role, got %d", recorderAdminRole.Code)
	}
}

func TestCommentsVisibilityAndEventsEndpoint(t *testing.T) {
	mux := setupTicketRouter()
	id := createTicketForHandlerTest(t, mux)

	publicComment := []byte(`{"author_id":"user-1","body":"Need update","is_internal":false}`)
	requestPublic := httptest.NewRequest(http.MethodPost, "/api/v1/tickets/"+id+"/comments", bytes.NewReader(publicComment))
	requestPublic.Header.Set("Content-Type", "application/json")
	recorderPublic := httptest.NewRecorder()
	mux.ServeHTTP(recorderPublic, requestPublic)
	if recorderPublic.Code != http.StatusCreated {
		t.Fatalf("expected 201 for public comment, got %d", recorderPublic.Code)
	}

	internalComment := []byte(`{"author_id":"agent-1","body":"Internal note","is_internal":true}`)
	requestInternalByClient := httptest.NewRequest(http.MethodPost, "/api/v1/tickets/"+id+"/comments", bytes.NewReader(internalComment))
	requestInternalByClient.Header.Set("Content-Type", "application/json")
	requestInternalByClient.Header.Set("X-User-Role", "client")
	recorderInternalByClient := httptest.NewRecorder()
	mux.ServeHTTP(recorderInternalByClient, requestInternalByClient)
	if recorderInternalByClient.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for client creating internal comment, got %d", recorderInternalByClient.Code)
	}

	requestInternalByAgent := httptest.NewRequest(http.MethodPost, "/api/v1/tickets/"+id+"/comments", bytes.NewReader(internalComment))
	requestInternalByAgent.Header.Set("Content-Type", "application/json")
	requestInternalByAgent.Header.Set("X-User-Role", "agent")
	recorderInternalByAgent := httptest.NewRecorder()
	mux.ServeHTTP(recorderInternalByAgent, requestInternalByAgent)
	if recorderInternalByAgent.Code != http.StatusCreated {
		t.Fatalf("expected 201 for agent internal comment, got %d", recorderInternalByAgent.Code)
	}

	requestListClient := httptest.NewRequest(http.MethodGet, "/api/v1/tickets/"+id+"/comments", nil)
	requestListClient.Header.Set("X-User-Role", "client")
	recorderListClient := httptest.NewRecorder()
	mux.ServeHTTP(recorderListClient, requestListClient)
	if recorderListClient.Code != http.StatusOK {
		t.Fatalf("expected 200 for client comment list, got %d", recorderListClient.Code)
	}
	var clientComments []ticket.Comment
	if err := json.Unmarshal(recorderListClient.Body.Bytes(), &clientComments); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if len(clientComments) != 1 {
		t.Fatalf("expected 1 visible comment for client, got %d", len(clientComments))
	}

	requestListAgent := httptest.NewRequest(http.MethodGet, "/api/v1/tickets/"+id+"/comments", nil)
	requestListAgent.Header.Set("X-User-Role", "agent")
	recorderListAgent := httptest.NewRecorder()
	mux.ServeHTTP(recorderListAgent, requestListAgent)
	if recorderListAgent.Code != http.StatusOK {
		t.Fatalf("expected 200 for agent comment list, got %d", recorderListAgent.Code)
	}
	var agentComments []ticket.Comment
	if err := json.Unmarshal(recorderListAgent.Body.Bytes(), &agentComments); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if len(agentComments) != 2 {
		t.Fatalf("expected 2 visible comments for agent, got %d", len(agentComments))
	}

	requestEvents := httptest.NewRequest(http.MethodGet, "/api/v1/tickets/"+id+"/events", nil)
	recorderEvents := httptest.NewRecorder()
	mux.ServeHTTP(recorderEvents, requestEvents)
	if recorderEvents.Code != http.StatusOK {
		t.Fatalf("expected 200 for events list, got %d", recorderEvents.Code)
	}
	var events []ticket.TicketEvent
	if err := json.Unmarshal(recorderEvents.Body.Bytes(), &events); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if len(events) == 0 {
		t.Fatal("expected non-empty event list")
	}
}
