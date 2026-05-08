-include .env
export

GOBIN := $(shell go env GOPATH)/bin

.PHONY: up down migrate-up migrate-down generate run dev test test-coverage

up:
	docker compose up -d

down:
	docker compose down

migrate-up:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="$(DATABASE_URL)" $(GOBIN)/goose -dir migrations up

migrate-down:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="$(DATABASE_URL)" $(GOBIN)/goose -dir migrations down

generate:
	$(GOBIN)/gqlgen generate
	$(GOBIN)/sqlc generate

run:
	go run ./cmd/server

dev:
	$(GOBIN)/gqlgen generate
	$(GOBIN)/sqlc generate
	go run ./cmd/server

test:
	go test ./... -v

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out
