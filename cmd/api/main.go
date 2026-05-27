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

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do/v2"

	"github.com/MayMustAI/core/internal/config"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	if err := run(); err != nil {
		slog.Error("api server failed", slog.Any("error", err))
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	injector := buildInjector(ctx)
	defer func() {
		if report := injector.Shutdown(); !report.Succeed {
			slog.Error("dependency shutdown failed", slog.Any("error", report))
		}
	}()

	// Resolve startup-critical dependencies eagerly so configuration errors or
	// unavailable backends fail fast, before the server accepts traffic.
	if _, err := do.Invoke[*config.Config](injector); err != nil {
		return err
	}
	pool, err := do.Invoke[*pgxpool.Pool](injector)
	if err != nil {
		return err
	}
	defer pool.Close()
	if _, err := do.Invoke[*oidc.Provider](injector); err != nil {
		return err
	}

	server, err := do.Invoke[*nethttp.Server](injector)
	if err != nil {
		return err
	}

	serverErr := make(chan error, 1)
	go func() {
		slog.Info("starting api server", slog.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, nethttp.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return server.Shutdown(shutdownCtx)
}
