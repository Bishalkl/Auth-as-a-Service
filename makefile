# Set app-related paths
APP_NAME = auth-service
CMD_PATH = ./cmd/server

# .env support
ENV_FILE = .env

# DB connection from .env
DB_DRIVER = postgres
DB_STRING = "host=$${DB_HOST} port=$${DB_PORT} user=$${DB_USER} password=$${DB_PASSWORD} dbname=$${DB_NAME} sslmode=disable"

# Goose migration directory
MIGRATIONS_DIR = ./migrations

# Load env
load-env:
	@export $(shell cat $(ENV_FILE) | xargs)

# Run the application
run:
	go run $(CMD_PATH)/main.go

# Build the application binary
build:
	go build -o bin/$(APP_NAME) $(CMD_PATH)/main.go

# Format the code
fmt:
	go fmt ./...

# Run tests with verbose output
test:
	go test -v ./...

# Run tests with coverage
test-cover:
	go test -v -cover ./...

# Lint code (assuming golangci-lint is installed)
lint:
	golangci-lint run ./...

# Goose: create a new SQL migration
migrate-create:
	goose -dir $(MIGRATIONS_DIR) create $(name) sql

# Goose: run migrations (up)
migrate-up:
	@export $(shell cat $(ENV_FILE) | xargs) && \
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_STRING) up

# Goose: rollback last migration (down)
migrate-down:
	@export $(shell cat $(ENV_FILE) | xargs) && \
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_STRING) down

# Goose: check status
migrate-status:
	@export $(shell cat $(ENV_FILE) | xargs) && \
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_STRING) status

# Goose: migrate up to a specific version
migrate-up-to:
	@export $(shell cat $(ENV_FILE) | xargs) && \
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_STRING) up $(version)

# Goose: rollback to a specific version
migrate-down-to:
	@export $(shell cat $(ENV_FILE) | xargs) && \
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_STRING) down $(version)

# Goose: reset DB by rolling all down then all up
migrate-reset:
	@export $(shell cat $(ENV_FILE) | xargs) && \
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_STRING) reset

# Goose: print current DB version
migrate-version:
	@export $(shell cat $(ENV_FILE) | xargs) && \
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_STRING) version

# Goose: fix duplicate migration numbering (optional)
migrate-fix:
	goose -dir $(MIGRATIONS_DIR) fix