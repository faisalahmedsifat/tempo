# Tempo Webhook Scheduler Makefile

.PHONY: build clean test install help

# Default target
all: build

# Build the CLI
build:
	@echo "Building Tempo CLI..."
	go build -o tempo cmd/cli/main.go
	@echo "✓ Built tempo CLI"

# Build the original test application
build-test:
	@echo "Building Tempo test application..."
	go build -o tempo-test cmd/main.go
	@echo "✓ Built tempo-test"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f tempo tempo-test
	@echo "✓ Cleaned build artifacts"

# Install CLI globally
install:
	@echo "Installing Tempo CLI globally..."
	go install ./cmd/cli
	@echo "✓ Installed tempo CLI"

# Run tests
test:
	@echo "Running tests..."
	go test ./...
	@echo "✓ Tests completed"

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run
	@echo "✓ Linting completed"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "✓ Code formatted"

# Show help
help:
	@echo "Tempo Webhook Scheduler - Available targets:"
	@echo "  build      - Build the CLI application"
	@echo "  build-test - Build the original test application"
	@echo "  clean      - Remove build artifacts"
	@echo "  install    - Install CLI globally"
	@echo "  test       - Run tests"
	@echo "  lint       - Run linter"
	@echo "  fmt        - Format code"
	@echo "  help       - Show this help message"

# Development helpers
dev-setup:
	@echo "Setting up development environment..."
	go mod tidy
	go mod download
	@echo "✓ Development environment ready"

# Quick test of CLI
test-cli: build
	@echo "Testing CLI functionality..."
	./tempo --help
	./tempo list
	@echo "✓ CLI test completed"