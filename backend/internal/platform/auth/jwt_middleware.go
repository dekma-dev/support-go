package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Claims struct {
	Subject string
	Role    string
}

type claimsContextKey struct{}

var errInvalidToken = errors.New("invalid token")

func NewJWTMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			rawAuth := strings.TrimSpace(request.Header.Get("Authorization"))
			if rawAuth == "" {
				next.ServeHTTP(writer, request)
				return
			}

			if !strings.HasPrefix(strings.ToLower(rawAuth), "bearer ") {
				writeUnauthorized(writer)
				return
			}

			token := strings.TrimSpace(rawAuth[len("Bearer "):])
			claims, err := parseAndValidateHS256(token, secret, time.Now().UTC())
			if err != nil {
				writeUnauthorized(writer)
				return
			}

			ctx := context.WithValue(request.Context(), claimsContextKey{}, claims)
			next.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}

func ClaimsFromContext(ctx context.Context) (Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey{}).(Claims)
	return claims, ok
}

func parseAndValidateHS256(token string, secret string, now time.Time) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, errInvalidToken
	}

	headerPayload := parts[0] + "." + parts[1]
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return Claims{}, errInvalidToken
	}

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(headerPayload))
	expected := mac.Sum(nil)
	if !hmac.Equal(signature, expected) {
		return Claims{}, errInvalidToken
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Claims{}, errInvalidToken
	}
	var header map[string]any
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return Claims{}, errInvalidToken
	}
	alg, _ := header["alg"].(string)
	if alg != "HS256" {
		return Claims{}, errInvalidToken
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Claims{}, errInvalidToken
	}
	var payload map[string]any
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return Claims{}, errInvalidToken
	}

	if exp, ok := numericClaim(payload["exp"]); ok && now.Unix() >= int64(exp) {
		return Claims{}, errInvalidToken
	}
	if nbf, ok := numericClaim(payload["nbf"]); ok && now.Unix() < int64(nbf) {
		return Claims{}, errInvalidToken
	}

	subject, _ := payload["sub"].(string)
	role, _ := payload["role"].(string)
	return Claims{
		Subject: strings.TrimSpace(subject),
		Role:    strings.ToLower(strings.TrimSpace(role)),
	}, nil
}

func numericClaim(value any) (int, bool) {
	switch parsed := value.(type) {
	case float64:
		return int(parsed), true
	case json.Number:
		intValue, err := parsed.Int64()
		if err != nil {
			return 0, false
		}
		return int(intValue), true
	case string:
		intValue, err := strconv.Atoi(strings.TrimSpace(parsed))
		if err != nil {
			return 0, false
		}
		return int(intValue), true
	default:
		return 0, false
	}
}

func writeUnauthorized(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(writer).Encode(map[string]string{"error": "unauthorized"})
}
