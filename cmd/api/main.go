package main

import (
	"context"
	"errors"
	"log/slog"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpserver "github.com/MayMustAI/core/internal/http"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	addr := ":8080"
	if value := os.Getenv("HTTP_ADDR"); value != "" {
		addr = value
	}

	server := &nethttp.Server{
		Addr:              addr,
		Handler:           httpserver.NewRouter(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		slog.Info("starting api server", slog.String("addr", addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, nethttp.ErrServerClosed) {
			slog.Error("api server failed", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("api server shutdown failed", slog.Any("error", err))
		os.Exit(1)
	}
}
