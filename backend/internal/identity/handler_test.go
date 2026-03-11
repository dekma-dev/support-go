package identity

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	platformauth "support-go/backend/internal/platform/auth"
)

const testJWTSecret = "test-jwt-secret"

func newTestHandler() *Handler {
	fixedNow := time.Date(2026, time.March, 11, 10, 0, 0, 0, time.UTC)
	return &Handler{
		secret:     testJWTSecret,
		now:        func() time.Time { return fixedNow },
		accessTTL:  15 * time.Minute,
		refreshTTL: 24 * time.Hour,
	}
}

func TestLoginIssuesAccessAndRefreshTokens(t *testing.T) {
	handler := newTestHandler()

	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(`{"subject":"agent-1","role":"agent"}`))
	recorder := httptest.NewRecorder()

	handler.login(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	var response tokenResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	accessClaims, err := platformauth.ParseAndValidateHS256(response.AccessToken, testJWTSecret, handler.now().Add(time.Minute))
	if err != nil {
		t.Fatalf("validate access token: %v", err)
	}
	if accessClaims.TokenType != "access" || accessClaims.Role != "agent" || accessClaims.Subject != "agent-1" {
		t.Fatalf("unexpected access claims: %+v", accessClaims)
	}

	refreshClaims, err := platformauth.ParseAndValidateHS256(response.RefreshToken, testJWTSecret, handler.now().Add(time.Minute))
	if err != nil {
		t.Fatalf("validate refresh token: %v", err)
	}
	if refreshClaims.TokenType != "refresh" || refreshClaims.Role != "agent" || refreshClaims.Subject != "agent-1" {
		t.Fatalf("unexpected refresh claims: %+v", refreshClaims)
	}
}

func TestLoginRejectsInvalidRole(t *testing.T) {
	handler := newTestHandler()

	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(`{"subject":"user-1","role":"owner"}`))
	recorder := httptest.NewRecorder()

	handler.login(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestRefreshRotatesTokens(t *testing.T) {
	handler := newTestHandler()
	loginResponse, err := handler.issueTokenPair("client-1", "client")
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
	loginResponse, err := handler.issueTokenPair("client-1", "client")
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
