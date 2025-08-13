.PHONY: help build run test clean deps

# Default target
help:
	@echo "Available commands:"
	@echo "  deps     - Install Go dependencies"
	@echo "  build    - Build the application"
	@echo "  run      - Run the application"
	@echo "  test     - Run tests"
	@echo "  clean    - Clean build artifacts"
	@echo "  dev      - Run in development mode with auto-reload"

# Install dependencies
deps:
	go mod tidy
	go mod download

# Build the application
build:
	go build -o dream-job-calculator .

# Run the application
run:
	go run .

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -f dream-job-calculator

# Development mode (requires air for hot reload)
dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Or run: make run"; \
	fi

# Database setup
db-setup:
	@echo "Setting up database..."
	@echo "1. Create database: CREATE DATABASE dream_job_calculator;"
	@echo "2. Run schema: psql -d dream_job_calculator -f schema.sql"
	@echo "3. Update .env file with your database credentials"

# Install development tools
install-tools:
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
