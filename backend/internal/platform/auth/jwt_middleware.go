package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type claimsContextKey struct{}

func isPublicRoute(path string) bool {
	if strings.HasPrefix(path, "/api/v1/auth/") {
		return true
	}
	switch path {
	case "/healthz", "/readyz":
		return true
	}
	return false
}

func NewJWTMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if isPublicRoute(request.URL.Path) {
				next.ServeHTTP(writer, request)
				return
			}

			rawAuth := strings.TrimSpace(request.Header.Get("Authorization"))
			if rawAuth == "" {
				writeUnauthorized(writer)
				return
			}

			if !strings.HasPrefix(strings.ToLower(rawAuth), "bearer ") {
				writeUnauthorized(writer)
				return
			}

			token := strings.TrimSpace(rawAuth[len("Bearer "):])
			claims, err := ParseAndValidateHS256(token, secret, time.Now().UTC())
			if err != nil || (claims.TokenType != "" && claims.TokenType != "access") {
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

func writeUnauthorized(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(writer).Encode(map[string]string{"error": "unauthorized"})
}
