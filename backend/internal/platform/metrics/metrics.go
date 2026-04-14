package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "support_go_http_requests_total",
			Help: "Total number of HTTP requests processed by the API, labeled by method, path, and status.",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "support_go_http_request_duration_seconds",
			Help:    "Histogram of HTTP request durations, labeled by method and path.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	TicketsCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "support_go_tickets_created_total",
			Help: "Total number of tickets created since process start.",
		},
	)

	TicketsByStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "support_go_tickets_by_status",
			Help: "Current number of tickets grouped by status.",
		},
		[]string{"status"},
	)
)

func Register() {
	prometheus.MustRegister(
		HTTPRequestsTotal,
		HTTPRequestDuration,
		TicketsCreatedTotal,
		TicketsByStatus,
	)
}

func Handler() http.Handler {
	return promhttp.Handler()
}
