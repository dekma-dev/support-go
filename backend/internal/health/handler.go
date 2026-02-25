package health

import (
	"encoding/json"
	"net/http"
	"time"
)

type healthResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Timestamp string `json:"timestamp"`
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", func(writer http.ResponseWriter, _ *http.Request) {
		writeHealth(writer, http.StatusOK)
	})

	mux.HandleFunc("/readyz", func(writer http.ResponseWriter, _ *http.Request) {
		writeHealth(writer, http.StatusOK)
	})
}

func writeHealth(writer http.ResponseWriter, statusCode int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)

	_ = json.NewEncoder(writer).Encode(healthResponse{
		Status:    "ok",
		Service:   "support-go-api",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

