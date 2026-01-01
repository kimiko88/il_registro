.PHONY: build run test lint docker-build docker-up migrate-up migrate-down

BINARY_NAME=api-server
DOCKER_IMAGE=registro-backend:latest

build:
	go build -o bin/$(BINARY_NAME) cmd/api-server/main.go

run:
	go run cmd/api-server/main.go


# Enhanced Testing
test:
	go test -v ./tests/... ./internal/...

test-unit:
	go test -v ./tests/unit/...

test-integration:
	go test -v ./tests/integration/...

test-coverage:
	go test -coverprofile=coverage.out ./internal/...
	go tool cover -html=coverage.out

lint:
	golangci-lint run

docker-build:
	docker build -f docker/Dockerfile.prod -t $(DOCKER_IMAGE) .

docker-up:
	DOCKER_BUILDKIT=0 COMPOSE_DOCKER_CLI_BUILD=0 docker compose -p registro-backend --env-file .env -f docker/docker-compose.yml up --build

# Migration commands (using golang-migrate)
# Install: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
MIGRATE_CMD=migrate -path migrations -database "postgres://postgres:password@localhost:5432/registro?sslmode=disable"

migrate-up:
	$(MIGRATE_CMD) up

migrate-down:
	$(MIGRATE_CMD) down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)
