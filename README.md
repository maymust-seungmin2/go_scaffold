# Go Scaffold 기술 스택 요약

```text
HTTP:          chi
Config:        spf13/viper
DB:            PostgreSQL + pgx
SQL Codegen:   sqlc
Migration:     golang-migrate
Logging:       log/slog + samber/slog-chi
Log pipeline:  slog-multi + slog-formatter + slog-sampling
Errors:        samber/oops
DI:            samber/do
Helpers:       samber/lo
Cache:         samber/hot
Functional:    samber/mo
Streams:       samber/ro
Validation:    go-playground/validator
Testing:       testing + testify + testcontainers-go
Lint:          golangci-lint
Security:      govulncheck
Docs:          swaggo/swag + http-swagger, optional
```

1. Go 의존성 설치
```bash
# HTTP 라우터: chi 본체와 기본 미들웨어 패키지
go get github.com/go-chi/chi/v5 github.com/go-chi/chi/v5/middleware

# PostgreSQL 드라이버/커넥션 풀: pgx, pgxpool
go get github.com/jackc/pgx/v5 github.com/jackc/pgx/v5/pgxpool

# 설정 관리: spf13/viper
go get github.com/spf13/viper

# 로깅: slog용 chi 미들웨어와 로그 파이프라인 구성 요소
go get github.com/samber/slog-chi github.com/samber/slog-multi github.com/samber/slog-formatter github.com/samber/slog-sampling

# 에러 처리: 구조화된 에러/스택/컨텍스트
go get github.com/samber/oops

# 의존성 주입: samber/do v2
go get github.com/samber/do/v2

# 헬퍼 함수: 컬렉션/함수형 유틸리티
go get github.com/samber/lo

# 캐시: 인메모리 캐시
go get github.com/samber/hot

# 함수형 안전 타입: Option/Result/Either 등
go get github.com/samber/mo

# 리액티브 스트림: 이벤트/스트림 파이프라인
go get github.com/samber/ro

# 검증: struct/request validation
go get github.com/go-playground/validator/v10

# 테스트: assertion/mock과 PostgreSQL testcontainers
go get github.com/stretchr/testify github.com/testcontainers/testcontainers-go github.com/testcontainers/testcontainers-go/modules/postgres

# Swagger/OpenAPI 런타임 UI: 문서 서버가 필요할 때 사용
go get github.com/swaggo/http-swagger github.com/swaggo/files

# go.mod/go.sum 정리
go mod tidy
```
