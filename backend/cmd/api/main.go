package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"support-go/backend/internal/health"
	"support-go/backend/internal/platform/config"
	platformhttp "support-go/backend/internal/platform/http"
	platformkafka "support-go/backend/internal/platform/kafka"
	"support-go/backend/internal/ticket"
	"support-go/backend/internal/ticket/postgres"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	logger.Info("starting support-go api", "port", cfg.HTTPPort, "env", cfg.Environment)

	if cfg.DatabaseURL == "" {
		logger.Error("DATABASE_URL is required")
		os.Exit(1)
	}

	dbPool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	if pingErr := dbPool.Ping(context.Background()); pingErr != nil {
		logger.Error("postgres ping failed", "error", pingErr)
		os.Exit(1)
	}
	logger.Info("postgres connected")

	var publisher ticket.EventPublisher
	var kafkaPublisher *platformkafka.Publisher
	brokers := platformkafka.ParseBrokers(cfg.KafkaBrokers)
	if len(brokers) == 0 {
		logger.Warn("kafka is not configured, domain events publishing disabled")
		publisher = platformkafka.NewNoopPublisher()
	} else {
		kafkaPublisher = platformkafka.NewPublisher(brokers)
		defer kafkaPublisher.Close()
		publisher = kafkaPublisher
		logger.Info("kafka publisher initialized", "brokers", brokers)
	}

	mux := http.NewServeMux()
	health.RegisterRoutes(mux)

	ticketRepository := postgres.NewRepository(dbPool)
	commentRepository := postgres.NewCommentRepository(dbPool)
	auditRepository := postgres.NewAuditRepository(dbPool)
	ticketService := ticket.NewServiceWithDependenciesAndPublisher(ticketRepository, commentRepository, auditRepository, publisher)
	ticket.RegisterRoutes(mux, ticketService)

	server := platformhttp.NewServer(cfg.HTTPPort, mux)
	serverErr := make(chan error, 1)

	go func() {
		serverErr <- server.ListenAndServe()
	}()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server crashed", "error", err)
			os.Exit(1)
		}
	case sig := <-shutdownSignal:
		logger.Info("shutdown signal received", "signal", sig.String())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("server stopped")
}
