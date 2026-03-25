package ticket

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	platformauth "support-go/backend/internal/platform/auth"
)

type Handler struct {
	service *Service
}

func RegisterRoutes(mux *http.ServeMux, service *Service) {
	handler := &Handler{service: service}
	mux.HandleFunc("/api/v1/tickets", handler.ticketsCollection)
	mux.HandleFunc("/api/v1/tickets/", handler.ticketsWithID)
}

type createTicketRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    Priority `json:"priority"`
	SLADueAt    *string  `json:"sla_due_at"`
}

type updateTicketRequest struct {
	Title       *string   `json:"title"`
	Description *string   `json:"description"`
	Priority    *Priority `json:"priority"`
}

type assignTicketRequest struct {
	AssigneeID string `json:"assignee_id"`
}

type statusTicketRequest struct {
	Status Status `json:"status"`
}

type createCommentRequest struct {
	Body       string `json:"body"`
	IsInternal bool   `json:"is_internal"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (handler *Handler) ticketsCollection(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		handler.createTicket(writer, request)
	case http.MethodGet:
		writeJSON(writer, http.StatusOK, handler.service.List())
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (handler *Handler) ticketsWithID(writer http.ResponseWriter, request *http.Request) {
	path := strings.Trim(request.URL.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		writeJSON(writer, http.StatusNotFound, errorResponse{Error: "route not found"})
		return
	}

	id := parts[3]
	if id == "" {
		writeValidationError(writer, "id is required")
		return
	}

	if len(parts) == 4 {
		handler.ticketByID(writer, request, id)
		return
	}

	subresource := parts[4]
	switch subresource {
	case "assign":
		if request.Method != http.MethodPatch {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handler.assignTicket(writer, request, id)
	case "status":
		if request.Method != http.MethodPatch {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handler.changeTicketStatus(writer, request, id)
	case "comments":
		switch request.Method {
		case http.MethodPost:
			handler.createComment(writer, request, id)
		case http.MethodGet:
			handler.listComments(writer, request, id)
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
		}
	case "events":
		if request.Method != http.MethodGet {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handler.listEvents(writer, request, id)
	default:
		writeJSON(writer, http.StatusNotFound, errorResponse{Error: "route not found"})
	}
}

func (handler *Handler) ticketByID(writer http.ResponseWriter, request *http.Request, id string) {
	switch request.Method {
	case http.MethodGet:
		ticketValue, err := handler.service.GetByID(id)
		if err != nil {
			handleServiceError(writer, err)
			return
		}
		writeJSON(writer, http.StatusOK, ticketValue)
	case http.MethodPatch:
		handler.patchTicket(writer, request, id)
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (handler *Handler) createTicket(writer http.ResponseWriter, request *http.Request) {
	claims, ok := platformauth.ClaimsFromContext(request.Context())
	if !ok || claims.Subject == "" {
		writeJSON(writer, http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
		return
	}

	var body createTicketRequest
	if err := decodeJSON(request, &body); err != nil {
		writeValidationError(writer, err.Error())
		return
	}

	var slaDueAt *time.Time
	if body.SLADueAt != nil {
		parsed, err := time.Parse(time.RFC3339, *body.SLADueAt)
		if err != nil {
			writeValidationError(writer, "sla_due_at must be RFC3339 timestamp")
			return
		}
		parsedUTC := parsed.UTC()
		slaDueAt = &parsedUTC
	}

	ticketValue, err := handler.service.Create(CreateInput{
		Title:       body.Title,
		Description: body.Description,
		Priority:    body.Priority,
		RequesterID: claims.Subject,
		SLADueAt:    slaDueAt,
	})
	if err != nil {
		handleServiceError(writer, err)
		return
	}

	writeJSON(writer, http.StatusCreated, ticketValue)
}

func (handler *Handler) patchTicket(writer http.ResponseWriter, request *http.Request, id string) {
	var body updateTicketRequest
	if err := decodeJSON(request, &body); err != nil {
		writeValidationError(writer, err.Error())
		return
	}

	ticketValue, err := handler.service.Update(id, UpdateInput{
		Title:       body.Title,
		Description: body.Description,
		Priority:    body.Priority,
	})
	if err != nil {
		handleServiceError(writer, err)
		return
	}

	writeJSON(writer, http.StatusOK, ticketValue)
}

func (handler *Handler) assignTicket(writer http.ResponseWriter, request *http.Request, id string) {
	claims, ok := platformauth.ClaimsFromContext(request.Context())
	if !ok || claims.Subject == "" {
		writeJSON(writer, http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
		return
	}
	if !canManageAssignmentsAndStatus(roleFromRequest(request)) {
		writeJSON(writer, http.StatusForbidden, errorResponse{Error: "forbidden: requires role agent or admin"})
		return
	}

	var body assignTicketRequest
	if err := decodeJSON(request, &body); err != nil {
		writeValidationError(writer, err.Error())
		return
	}

	ticketValue, err := handler.service.Assign(id, body.AssigneeID, claims.Subject)
	if err != nil {
		handleServiceError(writer, err)
		return
	}

	writeJSON(writer, http.StatusOK, ticketValue)
}

func (handler *Handler) changeTicketStatus(writer http.ResponseWriter, request *http.Request, id string) {
	claims, ok := platformauth.ClaimsFromContext(request.Context())
	if !ok || claims.Subject == "" {
		writeJSON(writer, http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
		return
	}
	if !canManageAssignmentsAndStatus(roleFromRequest(request)) {
		writeJSON(writer, http.StatusForbidden, errorResponse{Error: "forbidden: requires role agent or admin"})
		return
	}

	var body statusTicketRequest
	if err := decodeJSON(request, &body); err != nil {
		writeValidationError(writer, err.Error())
		return
	}

	ticketValue, err := handler.service.ChangeStatus(id, body.Status, claims.Subject)
	if err != nil {
		handleServiceError(writer, err)
		return
	}

	writeJSON(writer, http.StatusOK, ticketValue)
}

func (handler *Handler) createComment(writer http.ResponseWriter, request *http.Request, id string) {
	claims, ok := platformauth.ClaimsFromContext(request.Context())
	if !ok || claims.Subject == "" {
		writeJSON(writer, http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
		return
	}

	var body createCommentRequest
	if err := decodeJSON(request, &body); err != nil {
		writeValidationError(writer, err.Error())
		return
	}

	role := roleFromRequest(request)
	if body.IsInternal && !canManageAssignmentsAndStatus(role) {
		writeJSON(writer, http.StatusForbidden, errorResponse{Error: "forbidden: internal comments require role agent or admin"})
		return
	}

	commentValue, err := handler.service.AddComment(AddCommentInput{
		TicketID:   id,
		AuthorID:   claims.Subject,
		Body:       body.Body,
		IsInternal: body.IsInternal,
	})
	if err != nil {
		handleServiceError(writer, err)
		return
	}

	writeJSON(writer, http.StatusCreated, commentValue)
}

func (handler *Handler) listComments(writer http.ResponseWriter, request *http.Request, id string) {
	comments, err := handler.service.ListComments(id, roleFromRequest(request))
	if err != nil {
		handleServiceError(writer, err)
		return
	}

	writeJSON(writer, http.StatusOK, comments)
}

func (handler *Handler) listEvents(writer http.ResponseWriter, _ *http.Request, id string) {
	events, err := handler.service.ListEvents(id)
	if err != nil {
		handleServiceError(writer, err)
		return
	}

	writeJSON(writer, http.StatusOK, events)
}

func decodeJSON(request *http.Request, destination any) error {
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(destination)
}

func handleServiceError(writer http.ResponseWriter, err error) {
	if errors.Is(err, ErrNotFound) {
		writeJSON(writer, http.StatusNotFound, errorResponse{Error: err.Error()})
		return
	}
	if errors.Is(err, ErrValidation) {
		writeJSON(writer, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(writer, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
}

func writeValidationError(writer http.ResponseWriter, message string) {
	writeJSON(writer, http.StatusBadRequest, errorResponse{Error: message})
}

func writeJSON(writer http.ResponseWriter, statusCode int, payload any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	_ = json.NewEncoder(writer).Encode(payload)
}
