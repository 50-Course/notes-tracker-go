# We are using Mkaefile here to automate
# frequently used commands, and improve Developer experience
COMPOSE=docker compose
DB_CONTAINER=samba_db
SWAG = swag
GO = go

export $(shell sed 's/=.*//' config/.env)

run: 
	@echo "Starting services..."
	@$(COMPOSE) up --build -d

stop:
	@echo "Stopping services..."
	@$(COMPOSE) down

test:
	@echo "Running tests..."
	@go test ./... -v

build:
	@echo "Building binaries..."
	@go build -o bin/api ./api/gateway/server.go
	@go build -o bin/internal ./cmd/main.go

fmt:
	@echo "Formatting in progress..."
	@go fmt ./...

lint:
	@golangci-lint run ./...

swag:
	@echo "Generating Swagger docs..."
	@$(SWAG) init -g api/gateway/server.go --output ./docs

grpc:
	@echo "Starting gRPC server..."
	@go run cmd/main.go

gateway: swag
	@echo "Starting API Gateway..."
	@go run api/gateway/server.go

# think of this like django or mintlify build
serve: swag
	@echo "Starting gRPC Server & API Gateway..."
	@mkdir -p logs
	@nohup $(GO) run cmd/main.go > logs/grpc.log 2>&1 &  
	@nohup $(GO) run api/gateway/server.go > logs/gateway.log 2>&1 &
	@echo "Services started. Logs: logs/grpc.log, logs/gateway.log"

# stop both services (manual mode)
kill:
	@pkill -f "$(GO) run cmd/main.go" || true
	@pkill -f "$(GO) run api/gateway/server.go" || true
	@echo "gRPC & API Gateway stopped."

clean:
	@echo "Cleaning up..."
	# for now we are only cleaning logs, i wouldn't want to clean up bin only to find my bin/ folder cleanedup
	@rm -rf logs/*
	@echo "Cleaned complete"

# because make woun't stop using cached images, i want fresh installs
clean_build:
	@echo "Cleaning up and rebuilding image..."
	@$(COMPOSE) down -v
	@$(COMPOSE) build --no-cache

.PHONY: run stop test build fmt lint swag grpc gateway clean serve kill clean_build

