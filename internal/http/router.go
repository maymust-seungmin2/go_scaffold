// Package http provides the application's HTTP router and handlers.
package http

import (
	"log/slog"
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/samber/do/v2"
	slogchi "github.com/samber/slog-chi"
)

// NewRouter builds the application router. It resolves its dependencies from the
// injector so new handlers can pull services (config, DB pool, domain services)
// from the same container without changing this signature.
func NewRouter(i do.Injector) nethttp.Handler {
	logger := do.MustInvoke[*slog.Logger](i)

	r := chi.NewRouter()
	// RequestID runs before the logging middleware so the generated ID is in the
	// context and slog-chi can attach it to each request log line.
	r.Use(middleware.RequestID)
	r.Use(slogchi.NewWithConfig(logger, slogchi.Config{
		DefaultLevel:     slog.LevelInfo,
		ClientErrorLevel: slog.LevelWarn,
		ServerErrorLevel: slog.LevelError,
		WithRequestID:    true,
		// Health checks are high-frequency and low-signal; keep them out of logs.
		Filters: []slogchi.Filter{
			slogchi.IgnorePath("/healthz"),
		},
	}))
	r.Use(middleware.Recoverer)

	r.Get("/healthz", func(w nethttp.ResponseWriter, _ *nethttp.Request) {
		w.WriteHeader(nethttp.StatusOK)
		if _, err := w.Write([]byte("ok")); err != nil {
			logger.Warn("healthz write failed", slog.Any("error", err))
		}
	})

	return r
}
