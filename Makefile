.PHONY: run-api migrate-up migrate-down sqlc test lint tidy docker-up docker-down

run-api:
	go run ./cmd/api

migrate-up:
	migrate -path migrations -database "$${DATABASE_URL}" up

migrate-down:
	migrate -path migrations -database "$${DATABASE_URL}" down 1

sqlc:
	sqlc generate

test:
	go test ./...

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

docker-up:
	docker compose up -d

docker-down:
	docker compose down
