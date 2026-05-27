# Project Tech Stack Summary

```text
HTTP:          chi
Config:        spf13/viper + go-playground/validator
DB:            PostgreSQL + pgx
SQL Codegen:   sqlc
Migration:     golang-migrate
Logging:       log/slog + samber/slog-chi
Log pipeline:  slog-multi + slog-formatter + slog-sampling
Tracing:       OpenTelemetry
Errors:        samber/oops
DI:            samber/do
Helpers:       samber/lo
Cache:         samber/hot
Functional:    samber/mo
Streams:       samber/ro
RPC:           ConnectRPC + gRPC
API Schema:    Protobuf + Buf
Validation:    Protovalidate
Auth/IAM:      Keycloak + go-oidc
Testing:       testing + testify + testcontainers-go
Load Test:     vegeta
Lint:          golangci-lint
Security:      govulncheck
```

# Dependency Installation
```bash
# HTTP router: chi core package and standard middleware
go get github.com/go-chi/chi/v5 github.com/go-chi/chi/v5/middleware

# PostgreSQL driver and connection pool: pgx, pgxpool
go get github.com/jackc/pgx/v5 github.com/jackc/pgx/v5/pgxpool

# Configuration management: spf13/viper + validator for config validation
go get github.com/spf13/viper github.com/go-playground/validator/v10

# Logging: chi middleware for slog and log pipeline components
go get github.com/samber/slog-chi github.com/samber/slog-multi github.com/samber/slog-formatter github.com/samber/slog-sampling

# Observability: OpenTelemetry-based tracing, metrics, and logs instrumentation
go get \
  go.opentelemetry.io/otel \
  go.opentelemetry.io/otel/sdk \
  go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc \
  go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp \
  go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc \
  connectrpc.com/otelconnect

# Error handling: structured errors, stack traces, and context
go get github.com/samber/oops

# Dependency injection: samber/do v2
go get github.com/samber/do/v2

# Helper functions: collection and functional utilities
go get github.com/samber/lo

# Cache: in-memory cache
go get github.com/samber/hot

# Functional safety types: Option, Result, Either, and more
go get github.com/samber/mo

# Reactive streams: event and stream pipelines
go get github.com/samber/ro

# Protobuf / gRPC / ConnectRPC / Protovalidate
go get \
  google.golang.org/protobuf \
  google.golang.org/grpc \
  connectrpc.com/connect \
  github.com/bufbuild/protovalidate-go \
  connectrpc.com/validate

# Keycloak client library + OIDC/JWT verification
go get github.com/Nerzal/gocloak/v13 github.com/coreos/go-oidc/v3

# Testing: assertions, mocks, and PostgreSQL testcontainers
go get github.com/stretchr/testify github.com/testcontainers/testcontainers-go github.com/testcontainers/testcontainers-go/modules/postgres
```
