# ————————————————————————————————————————————————
# Set app-related paths
# ————————————————————————————————————————————————
APP_NAME = auth-service
CMD_PATH = ./cmd/server

# .env support
ENV_FILE = .env

# DB connection from .env
DB_DRIVER = postgres
DB_STRING = "host=$${DB_HOST} port=$${DB_PORT} user=$${DB_USER} password=$${DB_PASSWORD} dbname=$${DB_NAME} sslmode=disable"

# Goose migration directory
MIGRATIONS_DIR = ./migrations

# ————————————————————————————————————————————————
# 1️⃣ Environment
# ————————————————————————————————————————————————
# load-env:
#   Load your .env file into the shell so that
#   subsequent commands pick up DB_HOST, JWT_SECRET, etc.
load-env:
	@export $$(cat $(ENV_FILE) | xargs)

# ————————————————————————————————————————————————
# 2️⃣ Run & Build
# ————————————————————————————————————————————————
# run:       Run the server with live code
# build:     Build a production binary
# fmt:       go fmt all files
run:
	@echo "🚀 Starting server"
	go run $(CMD_PATH)/main.go

build:
	@echo "📦 Building binary"
	go build -o bin/$(APP_NAME) $(CMD_PATH)/main.go

fmt:
	@echo "🖌️  Formatting code"
	go fmt ./...

# ————————————————————————————————————————————————
# 3️⃣ Testing & Linting
# ————————————————————————————————————————————————
# test:        Run all tests verbosely
# test-cover:  Run tests with coverage
# lint:        Run golangci-lint
test:
	@echo "🧪 Running tests"
	go test -v ./...

test-cover:
	@echo "🧪 Running tests with coverage"
	go test -v -cover ./...

lint:
	@echo "🔍 Linting code"
	golangci-lint run ./...

# ————————————————————————————————————————————————
# 4️⃣ Goose Migrations
# ————————————————————————————————————————————————
# MIGRATIONS_DIR := ./migrations
# DB_DRIVER       := postgres
# DB_STRING       := "host=$${DB_HOST} port=$${DB_PORT} user=$${DB_USER} password=$${DB_PASSWORD} dbname=$${DB_NAME} sslmode=disable"

# migrate-create:  Create a new blank SQL migration
#   make migrate-create name=create_users_table
migrate-create:
	@echo "✏️  Creating new migration: $(name)"
	@goose -dir $(MIGRATIONS_DIR) create $(name) sql

# migrate-up:      Apply all “up” migrations
migrate-up:
	@echo "⬆️  Applying migrations up"
	@export $$(cat .env | xargs) && \
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_STRING) up

# migrate-down:    Roll back the last migration
migrate-down:
	@echo "⬇️ Rolling back last migration"
	@export $$(cat .env | xargs) && \
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_STRING) down

# migrate-status:  Show current migration status
migrate-status:
	@echo "📋 Migration status"
	@export $$(cat .env | xargs) && \
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_STRING) status

# migrate-up-to:   Apply migrations up to a specific version
#   make migrate-up-to version=20250424101431
migrate-up-to:
	@echo "⬆️  Migrating up to version $(version)"
	@export $$(cat .env | xargs) && \
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_STRING) up $(version)

# migrate-down-to: Roll back down to a specific version
#   make migrate-down-to version=20250424101431
migrate-down-to:
	@echo "⬇️ Rolling down to version $(version)"
	@export $$(cat .env | xargs) && \
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_STRING) down $(version)

# migrate-reset:   Roll ALL down, then ALL up (clears & rebuilds schema)
migrate-reset:
	@echo "⚠️  Resetting database (down all, up all)"
	@export $$(cat .env | xargs) && \
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_STRING) reset

# migrate-version: Print the current migration version
migrate-version:
	@echo "🔢 Current DB version"
	@export $$(cat .env | xargs) && \
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_STRING) version

# migrate-fix:     Renumber migrations if you ever duplicate
migrate-fix:
	@echo "🛠️  Fixing duplicate migration numbering"
	@goose -dir $(MIGRATIONS_DIR) fix
