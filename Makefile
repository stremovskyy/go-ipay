.PHONY: all test coverage lint clean build

# Default target
all: test lint build

# Build the application
build:
	go build -v ./...

# Run tests
test:
	go test -v -race ./...

# Run tests with coverage
coverage:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt -o coverage.html

# Run linter
lint:
	golangci-lint run

# Clean build artifacts
clean:
	go clean
	rm -f coverage.txt coverage.html

# Install development dependencies
setup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Update dependencies
deps:
	go mod tidy
	go mod verify

# Run security check
security:
	gosec ./...

# Generate documentation
docs:
	godoc -http=:6060

# Help target
help:
	@echo "Available targets:"
	@echo "  all       - Run tests, lint, and build"
	@echo "  build     - Build the application"
	@echo "  test      - Run tests"
	@echo "  coverage  - Run tests with coverage"
	@echo "  lint      - Run linter"
	@echo "  clean     - Clean build artifacts"
	@echo "  setup     - Install development dependencies"
	@echo "  deps      - Update dependencies"
	@echo "  security  - Run security check"
	@echo "  docs      - Generate documentation"
