# We are using Mkaefile here to automate
# frequently used commands, and improve Developer experience
COMPOSE=docker compose
DB_CONTAINER=samba_db

export $(shell sed 's/=.*//' config/.env)

run:
	@echo "Starting services..."
	@$(COMPOSE) up --build -d

stop:
	@echo "Stopping services..."
	@$(COMPOSE) down

test:
	@echo "Running tests..."
	@go test ./cmd/repository -v

build:
	@echo "Building binaries..."
	@go build -o bin/api ./api/gateway/server.go
	@go build -o bin/internal ./cmd/main.go

fmt:
	@echo "Formatting in progress..."
	@go fmt ./...

lint:
	@golangci-lint run ./...

.PHONY: run stop test build fmt lint

