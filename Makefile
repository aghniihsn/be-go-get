# Cinema Booking API Makefile

.PHONY: run build docs clean test generate-secret

# Default target
all: docs run

# Run the application
run:
	@echo "Starting Cinema Booking API..."
	go run main.go

# Build the application
build:
	@echo "Building Cinema Booking API..."
	go build -o bin/cinema-api main.go

# Generate Swagger documentation
docs:
	@echo "Generating Swagger documentation..."
	swag init
	@echo "Swagger docs generated! Access at: http://localhost:3000/swagger/"

# Generate JWT Secret
generate-secret:
	@echo "🔐 Generating secure JWT secrets..."
	@go run scripts/generate-secret/main.go

# Update JWT Secret in .env file
update-secret:
	@echo "🔄 Updating JWT secret in .env file..."
	@go run scripts/update-secret/main.go

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf bin/
	rm -rf docs/*.go docs/*.json docs/*.yaml

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Setup development environment
setup:
	@echo "Setting up development environment..."
	go mod init go-get-backend || true
	go get github.com/gofiber/fiber/v2
	go get github.com/joho/godotenv
	go get go.mongodb.org/mongo-driver/mongo
	go get github.com/swaggo/swag/cmd/swag
	go get github.com/swaggo/fiber-swagger
	go mod tidy

# Run in development mode with auto-reload (requires air)
dev:
	@echo "Starting development server..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not installed. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Or run with: make run"; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install from: https://golangci-lint.run/usage/install/"; \
	fi

# Show help
help:
	@echo "Available commands:"
	@echo "  run     - Run the application"
	@echo "  build   - Build the application"
	@echo "  docs    - Generate Swagger documentation"
	@echo "  clean   - Clean build artifacts"
	@echo "  deps    - Install dependencies"
	@echo "  test    - Run tests"
	@echo "  setup   - Setup development environment"
	@echo "  dev     - Run in development mode with auto-reload"
	@echo "  fmt     - Format code"
	@echo "  lint    - Lint code"
	@echo "  help    - Show this help message"
