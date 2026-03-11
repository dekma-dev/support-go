package identity

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	platformauth "support-go/backend/internal/platform/auth"
)

const (
	defaultAccessTTL  = 15 * time.Minute
	defaultRefreshTTL = 24 * time.Hour
)

type Handler struct {
	secret     string
	now        func() time.Time
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type loginRequest struct {
	Subject string `json:"subject"`
	Role    string `json:"role"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type tokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	ExpiresAt        string `json:"expires_at"`
	RefreshExpiresAt string `json:"refresh_expires_at"`
	Subject          string `json:"subject"`
	Role             string `json:"role"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func RegisterRoutes(mux *http.ServeMux, secret string) {
	handler := &Handler{
		secret:     secret,
		now:        func() time.Time { return time.Now().UTC() },
		accessTTL:  defaultAccessTTL,
		refreshTTL: defaultRefreshTTL,
	}

	mux.HandleFunc("/api/v1/auth/login", handler.login)
	mux.HandleFunc("/api/v1/auth/refresh", handler.refresh)
}

func (handler *Handler) login(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var input loginRequest
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeJSON(writer, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	subject := strings.TrimSpace(input.Subject)
	role := normalizeRole(input.Role)
	if subject == "" {
		writeJSON(writer, http.StatusBadRequest, errorResponse{Error: "subject is required"})
		return
	}
	if !isAllowedRole(role) {
		writeJSON(writer, http.StatusBadRequest, errorResponse{Error: "role must be one of client, agent, admin"})
		return
	}

	response, err := handler.issueTokenPair(subject, role)
	if err != nil {
		writeJSON(writer, http.StatusInternalServerError, errorResponse{Error: "failed to issue token"})
		return
	}

	writeJSON(writer, http.StatusOK, response)
}

func (handler *Handler) refresh(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var input refreshRequest
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeJSON(writer, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	now := handler.now()
	claims, err := platformauth.ParseAndValidateHS256(strings.TrimSpace(input.RefreshToken), handler.secret, now)
	if err != nil || claims.TokenType != "refresh" || claims.Subject == "" || !isAllowedRole(claims.Role) {
		writeJSON(writer, http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
		return
	}

	response, err := handler.issueTokenPair(claims.Subject, claims.Role)
	if err != nil {
		writeJSON(writer, http.StatusInternalServerError, errorResponse{Error: "failed to issue token"})
		return
	}

	writeJSON(writer, http.StatusOK, response)
}

func (handler *Handler) issueTokenPair(subject string, role string) (tokenResponse, error) {
	now := handler.now()
	accessExpiresAt := now.Add(handler.accessTTL)
	refreshExpiresAt := now.Add(handler.refreshTTL)

	accessToken, err := platformauth.IssueHS256Token(handler.secret, platformauth.Claims{
		Subject:   subject,
		Role:      role,
		TokenType: "access",
	}, now, accessExpiresAt)
	if err != nil {
		return tokenResponse{}, err
	}

	refreshToken, err := platformauth.IssueHS256Token(handler.secret, platformauth.Claims{
		Subject:   subject,
		Role:      role,
		TokenType: "refresh",
	}, now, refreshExpiresAt)
	if err != nil {
		return tokenResponse{}, err
	}

	return tokenResponse{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		TokenType:        "Bearer",
		ExpiresAt:        accessExpiresAt.Format(time.RFC3339),
		RefreshExpiresAt: refreshExpiresAt.Format(time.RFC3339),
		Subject:          subject,
		Role:             role,
	}, nil
}

func normalizeRole(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func isAllowedRole(role string) bool {
	switch role {
	case "client", "agent", "admin":
		return true
	default:
		return false
	}
}

func writeJSON(writer http.ResponseWriter, status int, payload any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(payload)
}
