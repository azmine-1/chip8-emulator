# CHIP-8 Emulator Makefile

# Binary name
BINARY_NAME=chip8-emulator
BINARY_WINDOWS=$(BINARY_NAME).exe
BINARY_UNIX=$(BINARY_NAME)

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-s -w"

# Default target
all: build

# Build the application
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_WINDOWS) .

# Build for Windows
build-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_WINDOWS) .

# Build for Linux
build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_UNIX) .

# Build for macOS
build-mac:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_UNIX) .

# Run the application
run:
	$(GOBUILD) -o $(BINARY_WINDOWS) .
	./$(BINARY_WINDOWS)

# Run with race detection
run-race:
	$(GOBUILD) -race -o $(BINARY_WINDOWS) .
	./$(BINARY_WINDOWS)

# Test the application
test:
	$(GOTEST) -v ./...

# Test with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_WINDOWS)
	rm -f $(BINARY_UNIX)
	rm -f coverage.out

# Install dependencies
deps:
	$(GOGET) -v -t -d ./...

# Tidy go modules
tidy:
	$(GOMOD) tidy

# Format code
fmt:
	$(GOCMD) fmt ./...

# Vet code
vet:
	$(GOCMD) vet ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Install golangci-lint (if not installed)
install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Development target - build and run
dev: build run

# Help target
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  build-windows  - Build for Windows"
	@echo "  build-linux    - Build for Linux"
	@echo "  build-mac      - Build for macOS"
	@echo "  run            - Build and run the application"
	@echo "  run-race       - Run with race detection"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  clean          - Clean build artifacts"
	@echo "  deps           - Install dependencies"
	@echo "  tidy           - Tidy go modules"
	@echo "  fmt            - Format code"
	@echo "  vet            - Vet code"
	@echo "  lint           - Lint code (requires golangci-lint)"
	@echo "  install-lint   - Install golangci-lint"
	@echo "  dev            - Build and run (development)"
	@echo "  help           - Show this help message"

.PHONY: all build build-windows build-linux build-mac run run-race test test-coverage clean deps tidy fmt vet lint install-lint dev help 