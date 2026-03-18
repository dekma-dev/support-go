package identity

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	platformauth "support-go/backend/internal/platform/auth"
)

const (
	defaultAccessTTL  = 15 * time.Minute
	defaultRefreshTTL = 24 * time.Hour
	userStatusActive  = "active"
)

type Handler struct {
	repo       Repository
	secret     string
	now        func() time.Time
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
	UserID           string `json:"user_id"`
	Email            string `json:"email"`
	Role             string `json:"role"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func RegisterRoutes(mux *http.ServeMux, repo Repository, secret string) {
	handler := &Handler{
		repo:       repo,
		secret:     secret,
		now:        func() time.Time { return time.Now().UTC() },
		accessTTL:  defaultAccessTTL,
		refreshTTL: defaultRefreshTTL,
	}

	mux.HandleFunc("/api/v1/auth/register", handler.register)
	mux.HandleFunc("/api/v1/auth/login", handler.login)
	mux.HandleFunc("/api/v1/auth/refresh", handler.refresh)
}

func (handler *Handler) register(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var input registerRequest
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeJSON(writer, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	email, err := normalizeEmail(input.Email)
	if err != nil {
		writeJSON(writer, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}
	if err := validatePassword(input.Password); err != nil {
		writeJSON(writer, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	role := normalizeRole(input.Role)
	if !isAllowedRole(role) {
		writeJSON(writer, http.StatusBadRequest, errorResponse{Error: "role must be one of client, agent, admin"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		writeJSON(writer, http.StatusInternalServerError, errorResponse{Error: "failed to hash password"})
		return
	}

	now := handler.now()
	user, err := handler.repo.CreateUser(request.Context(), User{
		ID:           newUserID(),
		Email:        email,
		PasswordHash: string(passwordHash),
		Role:         role,
		Status:       userStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if errors.Is(err, ErrEmailAlreadyExists) {
		writeJSON(writer, http.StatusConflict, errorResponse{Error: "email already exists"})
		return
	}
	if err != nil {
		writeJSON(writer, http.StatusInternalServerError, errorResponse{Error: "failed to create user"})
		return
	}

	response, err := handler.issueTokenPair(user)
	if err != nil {
		writeJSON(writer, http.StatusInternalServerError, errorResponse{Error: "failed to issue token"})
		return
	}

	writeJSON(writer, http.StatusCreated, response)
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

	email, err := normalizeEmail(input.Email)
	if err != nil {
		writeJSON(writer, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}
	if strings.TrimSpace(input.Password) == "" {
		writeJSON(writer, http.StatusBadRequest, errorResponse{Error: "password is required"})
		return
	}

	user, err := handler.repo.GetUserByEmail(request.Context(), email)
	if err != nil || !isActiveUser(user) {
		writeJSON(writer, http.StatusUnauthorized, errorResponse{Error: "invalid credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		writeJSON(writer, http.StatusUnauthorized, errorResponse{Error: "invalid credentials"})
		return
	}

	response, err := handler.issueTokenPair(user)
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
	if err != nil || claims.TokenType != "refresh" || claims.Subject == "" {
		writeJSON(writer, http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
		return
	}

	user, err := handler.repo.GetUserByID(request.Context(), claims.Subject)
	if err != nil || !isActiveUser(user) {
		writeJSON(writer, http.StatusUnauthorized, errorResponse{Error: "unauthorized"})
		return
	}

	response, err := handler.issueTokenPair(user)
	if err != nil {
		writeJSON(writer, http.StatusInternalServerError, errorResponse{Error: "failed to issue token"})
		return
	}

	writeJSON(writer, http.StatusOK, response)
}

func (handler *Handler) issueTokenPair(user User) (tokenResponse, error) {
	now := handler.now()
	accessExpiresAt := now.Add(handler.accessTTL)
	refreshExpiresAt := now.Add(handler.refreshTTL)

	accessToken, err := platformauth.IssueHS256Token(handler.secret, platformauth.Claims{
		Subject:   user.ID,
		Role:      user.Role,
		TokenType: "access",
	}, now, accessExpiresAt)
	if err != nil {
		return tokenResponse{}, err
	}

	refreshToken, err := platformauth.IssueHS256Token(handler.secret, platformauth.Claims{
		Subject:   user.ID,
		Role:      user.Role,
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
		UserID:           user.ID,
		Email:            user.Email,
		Role:             user.Role,
	}, nil
}

func normalizeEmail(value string) (string, error) {
	email := strings.ToLower(strings.TrimSpace(value))
	if email == "" {
		return "", errors.New("email is required")
	}
	parsed, err := mail.ParseAddress(email)
	if err != nil || !strings.EqualFold(strings.TrimSpace(parsed.Address), email) {
		return "", errors.New("email is invalid")
	}
	return email, nil
}

func validatePassword(value string) error {
	if len(strings.TrimSpace(value)) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
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

func isActiveUser(user User) bool {
	return strings.EqualFold(strings.TrimSpace(user.Status), userStatusActive)
}

func newUserID() string {
	bytes := make([]byte, 12)
	if _, err := rand.Read(bytes); err != nil {
		return "user_" + hex.EncodeToString([]byte(time.Now().UTC().Format("20060102150405")))
	}
	return "user_" + hex.EncodeToString(bytes)
}

func writeJSON(writer http.ResponseWriter, status int, payload any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(payload)
}

type memoryRepository struct {
	usersByID    map[string]User
	usersByEmail map[string]User
}

func newMemoryRepository() *memoryRepository {
	return &memoryRepository{
		usersByID:    make(map[string]User),
		usersByEmail: make(map[string]User),
	}
}

func (repository *memoryRepository) CreateUser(_ context.Context, user User) (User, error) {
	if _, exists := repository.usersByEmail[user.Email]; exists {
		return User{}, ErrEmailAlreadyExists
	}
	repository.usersByID[user.ID] = user
	repository.usersByEmail[user.Email] = user
	return user, nil
}

func (repository *memoryRepository) GetUserByEmail(_ context.Context, email string) (User, error) {
	user, ok := repository.usersByEmail[email]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return user, nil
}

func (repository *memoryRepository) GetUserByID(_ context.Context, id string) (User, error) {
	user, ok := repository.usersByID[id]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return user, nil
}
