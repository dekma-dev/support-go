package http

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"support-go/backend/internal/platform/metrics"
)

func NewMetricsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(recorder, r)

			path := normalizePath(r.URL.Path)
			duration := time.Since(start).Seconds()

			metrics.HTTPRequestsTotal.WithLabelValues(r.Method, path, strconv.Itoa(recorder.status)).Inc()
			metrics.HTTPRequestDuration.WithLabelValues(r.Method, path).Observe(duration)
		})
	}
}

// normalizePath collapses dynamic path segments (IDs) into a template
// so Prometheus labels stay bounded. Extend as new routes appear.
func normalizePath(path string) string {
	if path == "" {
		return "/"
	}

	segments := strings.Split(strings.Trim(path, "/"), "/")
	if len(segments) >= 4 && segments[0] == "api" && segments[1] == "v1" && segments[2] == "tickets" {
		// /api/v1/tickets/{id} or /api/v1/tickets/{id}/{sub}
		result := "/api/v1/tickets/{id}"
		if len(segments) >= 5 {
			result += "/" + segments[4]
		}
		return result
	}

	return path
}
