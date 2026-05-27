# AGENTS.md

Every agent writing code in this repository must always use the tech stack defined below.
Do not substitute arbitrary libraries/frameworks; if a new dependency is needed, confirm with the user first.

## Tech Stack

| Area | Technology |
| --- | --- |
| HTTP | chi (`github.com/go-chi/chi/v5`) |
| Config | spf13/viper + go-playground/validator |
| DB | PostgreSQL + pgx (`github.com/jackc/pgx/v5`) |
| SQL Codegen | sqlc |
| Migration | golang-migrate |
| Retry | cenkalti/backoff (`github.com/cenkalti/backoff/v4`) |
| Logging | `log/slog` + samber/slog-chi |
| Log pipeline | slog-multi + slog-formatter + slog-sampling |
| Tracing | OpenTelemetry |
| Errors | samber/oops |
| DI | samber/do v2 (`github.com/samber/do/v2`) |
| Helpers | samber/lo |
| Cache | samber/hot |
| Functional | samber/mo |
| Streams | samber/ro |
| RPC | ConnectRPC + gRPC |
| API Schema | Protobuf + Buf |
| Validation | Protovalidate |
| Auth/IAM | Keycloak + go-oidc |
| Testing | `testing` + testify + testcontainers-go |
| Load Test | vegeta |
| Task Runner | Taskfile |
| Lint | golangci-lint |
| Security | govulncheck |
