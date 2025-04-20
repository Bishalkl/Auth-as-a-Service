APP_NAME = auth-service
CMD_PATH = ./cmd/server


#Run the application
run:
	go run $(CMD_PATH)/main.go

# Build the application binary
build:
	go build -o bin/$(APP_NAME) $(CMD_PATH)/main.go

#Format the code 
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



# # Run DB migrations (requires golang-migrate installed)
# migrate-up:
# 	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

# migrate-down:
# 	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down

# migrate-force:
# 	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force 1

# # Generate mocks
# generate-mocks:
# 	mockgen -source=internal/services/auth_service.go -destination=mocks/mock_auth_service.go -package=mocks

# .PHONY: run build fmt test test-cover lint migrate-up migrate-down migrate-force generate-mocks

