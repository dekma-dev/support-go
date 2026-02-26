package ticket_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"support-go/backend/internal/ticket"
	"support-go/backend/internal/ticket/memory"
)

func setupTicketRouter() *http.ServeMux {
	repository := memory.NewRepository()
	service := ticket.NewService(repository)
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
