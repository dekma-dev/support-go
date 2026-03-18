package identity

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	platformauth "support-go/backend/internal/platform/auth"
)

const testJWTSecret = "test-jwt-secret"

func newTestHandler() *Handler {
	fixedNow := time.Date(2026, time.March, 11, 10, 0, 0, 0, time.UTC)
	return &Handler{
		repo:       newMemoryRepository(),
		secret:     testJWTSecret,
		now:        func() time.Time { return fixedNow },
		accessTTL:  15 * time.Minute,
		refreshTTL: 24 * time.Hour,
	}
}

func TestRegisterCreatesUserAndIssuesTokens(t *testing.T) {
	handler := newTestHandler()

	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(`{"email":"agent@example.com","password":"password123","role":"agent"}`))
	recorder := httptest.NewRecorder()

	handler.register(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", recorder.Code)
	}

	var response tokenResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Email != "agent@example.com" || response.Role != "agent" || response.UserID == "" {
		t.Fatalf("unexpected register response: %+v", response)
	}

	accessClaims, err := platformauth.ParseAndValidateHS256(response.AccessToken, testJWTSecret, handler.now().Add(time.Minute))
	if err != nil {
		t.Fatalf("validate access token: %v", err)
	}
	if accessClaims.TokenType != "access" || accessClaims.Role != "agent" || accessClaims.Subject != response.UserID {
		t.Fatalf("unexpected access claims: %+v", accessClaims)
	}
}

func TestRegisterRejectsDuplicateEmail(t *testing.T) {
	handler := newTestHandler()

	first := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(`{"email":"client@example.com","password":"password123","role":"client"}`))
	firstRecorder := httptest.NewRecorder()
	handler.register(firstRecorder, first)

	second := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(`{"email":"client@example.com","password":"password123","role":"client"}`))
	secondRecorder := httptest.NewRecorder()
	handler.register(secondRecorder, second)

	if secondRecorder.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", secondRecorder.Code)
	}
}

func TestLoginReturnsTokensForStoredUser(t *testing.T) {
	handler := newTestHandler()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	_, err = handler.repo.CreateUser(context.Background(), User{
		ID:           "user_123",
		Email:        "agent@example.com",
		PasswordHash: string(passwordHash),
		Role:         "agent",
		Status:       userStatusActive,
		CreatedAt:    handler.now(),
		UpdatedAt:    handler.now(),
	})
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(`{"email":"agent@example.com","password":"password123"}`))
	recorder := httptest.NewRecorder()

	handler.login(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
}

func TestRefreshRotatesTokens(t *testing.T) {
	handler := newTestHandler()
	user := User{
		ID:           "user_123",
		Email:        "client@example.com",
		PasswordHash: "ignored",
		Role:         "client",
		Status:       userStatusActive,
		CreatedAt:    handler.now(),
		UpdatedAt:    handler.now(),
	}
	_, err := handler.repo.CreateUser(context.Background(), user)
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}
	loginResponse, err := handler.issueTokenPair(user)
	if err != nil {
		t.Fatalf("issue token pair: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewBufferString(`{"refresh_token":"`+loginResponse.RefreshToken+`"}`))
	recorder := httptest.NewRecorder()

	handler.refresh(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	var response tokenResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.RefreshToken == loginResponse.RefreshToken {
		t.Fatal("expected refresh token rotation")
	}
}

func TestRefreshRejectsAccessToken(t *testing.T) {
	handler := newTestHandler()
	user := User{
		ID:           "user_123",
		Email:        "client@example.com",
		PasswordHash: "ignored",
		Role:         "client",
		Status:       userStatusActive,
		CreatedAt:    handler.now(),
		UpdatedAt:    handler.now(),
	}
	_, err := handler.repo.CreateUser(context.Background(), user)
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}
	loginResponse, err := handler.issueTokenPair(user)
	if err != nil {
		t.Fatalf("issue token pair: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewBufferString(`{"refresh_token":"`+loginResponse.AccessToken+`"}`))
	recorder := httptest.NewRecorder()

	handler.refresh(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
}
