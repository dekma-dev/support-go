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

	"support-go/backend/internal/health"
	"support-go/backend/internal/platform/config"
	platformhttp "support-go/backend/internal/platform/http"
	"support-go/backend/internal/ticket"
	"support-go/backend/internal/ticket/memory"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	logger.Info("starting support-go api", "port", cfg.HTTPPort, "env", cfg.Environment)

	mux := http.NewServeMux()
	health.RegisterRoutes(mux)

	ticketRepository := memory.NewRepository()
	ticketService := ticket.NewService(ticketRepository)
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
