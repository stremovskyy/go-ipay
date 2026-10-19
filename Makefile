.PHONY: all test coverage lint clean build setup deps security vuln docs help

# Default target
all: test lint build

# Build the application
build:
	go build -v ./...
	go -C examples/status_recorder build -v ./...

# Run tests
test:
	go test -v -race ./...
	go -C examples/status_recorder test -v -race ./...

# Run tests with coverage
coverage:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt -o coverage.html
	go -C examples/status_recorder test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# Run linter
lint:
	golangci-lint run

# Clean build artifacts
clean:
	go clean
	go -C examples/status_recorder clean
	rm -f coverage.txt coverage.html
	rm -f examples/status_recorder/coverage.txt

# Install development dependencies
setup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8

# Update dependencies
deps:
	go mod tidy
	go mod verify
	go -C examples/status_recorder mod tidy
	go -C examples/status_recorder mod verify

# Run security check
security:
	gosec -exclude-dir=private ./...
	cd examples/status_recorder && gosec ./...

# Run vulnerability checks with a patched Go toolchain.
vuln:
	govulncheck ./...
	cd examples/status_recorder && govulncheck ./...

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
	@echo "  vuln      - Check SDK and example dependency vulnerabilities"
	@echo "  docs      - Generate documentation"
