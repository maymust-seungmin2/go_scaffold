# Project Tech Stack Summary

```text
HTTP:          chi
Config:        spf13/viper + go-playground/validator
DB:            PostgreSQL + pgx
SQL Codegen:   sqlc
Migration:     golang-migrate
Retry:         cenkalti/backoff
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
Task Runner:   Taskfile
Lint:          golangci-lint
Security:      govulncheck
```
