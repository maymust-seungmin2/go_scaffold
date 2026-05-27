package main

import (
	"context"
	"fmt"
	"log/slog"
	nethttp "net/http"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do/v2"
	"github.com/samber/oops"

	"github.com/MayMustAI/core/internal/config"
	httpserver "github.com/MayMustAI/core/internal/http"
	"github.com/MayMustAI/core/internal/store"
)

// httpServerReadHeaderTimeout bounds how long the server waits for request
// headers, mitigating slow-client (Slowloris) resource exhaustion.
const httpServerReadHeaderTimeout = 5 * time.Second

// buildInjector registers every application service in the DI container.
// As the application grows (repositories, services, handlers, external
// clients), register them here; their dependencies are wired automatically by
// type inference and instantiated lazily on first use.
//
// ctx is the application's root context; providers use it so that long-running
// startup work (dependency retries) is cancelled on shutdown signals.
func buildInjector(ctx context.Context) *do.RootScope {
	injector := do.New()
	do.ProvideValue(injector, slog.Default())
	do.Provide(injector, provideConfig)
	do.Provide(injector, func(i do.Injector) (*pgxpool.Pool, error) {
		return provideDatabase(ctx, i)
	})
	do.Provide(injector, func(i do.Injector) (*oidc.Provider, error) {
		return provideOIDCProvider(ctx, i)
	})
	do.Provide(injector, provideRouter)
	do.Provide(injector, provideHTTPServer)
	return injector
}

// provideRouter builds the HTTP handler, passing the injector through so route
// handlers can resolve their own dependencies from the container.
func provideRouter(i do.Injector) (nethttp.Handler, error) {
	return httpserver.NewRouter(i), nil
}

// provideHTTPServer assembles the HTTP server from the resolved config and
// router. Lifecycle (start/shutdown) is owned by the caller at the composition
// root, not by the container.
func provideHTTPServer(i do.Injector) (*nethttp.Server, error) {
	cfg := do.MustInvoke[*config.Config](i)
	handler := do.MustInvoke[nethttp.Handler](i)

	return &nethttp.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: httpServerReadHeaderTimeout,
	}, nil
}

func provideConfig(_ do.Injector) (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, oops.In("startup").Wrapf(err, "load config")
	}
	return cfg, nil
}

// provideDatabase opens the connection pool and verifies connectivity with a
// retrying ping, so an unavailable database fails fast at startup.
func provideDatabase(ctx context.Context, i do.Injector) (*pgxpool.Pool, error) {
	cfg := do.MustInvoke[*config.Config](i)

	pool, err := store.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, oops.In("startup").Tags("database").Wrapf(err, "open database")
	}

	if err := retry(ctx, "ping database", func(ctx context.Context) error {
		pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		return pool.Ping(pingCtx)
	}); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

// provideOIDCProvider performs OIDC discovery against the Keycloak realm,
// retrying until it is reachable. The returned provider is reused for token
// verification. The issuer follows Keycloak's {KEYCLOAK_URL}/realms/{realm}.
func provideOIDCProvider(ctx context.Context, i do.Injector) (*oidc.Provider, error) {
	cfg := do.MustInvoke[*config.Config](i)
	issuer := fmt.Sprintf("%s/realms/%s", strings.TrimRight(cfg.KeycloakURL, "/"), cfg.KeycloakRealm)

	var provider *oidc.Provider
	if err := retry(ctx, "check keycloak", func(ctx context.Context) error {
		reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		p, err := oidc.NewProvider(reqCtx, issuer)
		if err != nil {
			return err
		}
		provider = p
		return nil
	}); err != nil {
		return nil, err
	}

	return provider, nil
}

// startupRetryMaxElapsed caps how long startup dependency checks (database,
// Keycloak) keep retrying with exponential backoff before failing fast. It
// tolerates dependencies that come up slightly later than the API process.
const startupRetryMaxElapsed = 30 * time.Second

// retry runs fn with exponential backoff until it succeeds or
// startupRetryMaxElapsed passes. It aborts early if ctx is cancelled, logs each
// failed attempt, and wraps the final error with name for context.
func retry(ctx context.Context, name string, fn func(context.Context) error) error {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = startupRetryMaxElapsed

	operation := func() error { return fn(ctx) }
	notify := func(err error, next time.Duration) {
		slog.Warn("startup check failed, retrying",
			slog.String("check", name),
			slog.Duration("retry_in", next),
			slog.Any("error", err))
	}

	if err := backoff.RetryNotify(operation, backoff.WithContext(b, ctx), notify); err != nil {
		return oops.
			In("startup").
			Code("dependency_unavailable").
			With("check", name).
			Wrapf(err, "startup dependency check failed")
	}
	return nil
}
