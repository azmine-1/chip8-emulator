# CHIP-8 Emulator Makefile

# Detect OS for better cross-platform support
ifeq ($(OS),Windows_NT)
    # Windows
    BINARY_NAME=chip8-emulator.exe
    RM=del /Q
    RMDIR=rmdir /S /Q
else
    # Unix-like systems (Linux, macOS, WSL)
    BINARY_NAME=chip8-emulator
    RM=rm -f
    RMDIR=rm -rf
endif

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-s -w"

# Default target - builds for current platform
all: build

# Build the application for current platform
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) .

# Build for Windows (from WSL)
build-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o chip8-emulator.exe .

# Build for Linux
build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o chip8-emulator .

# Build for macOS
build-mac:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o chip8-emulator .

# Build for ARM64 Linux (useful for some WSL setups)
build-linux-arm64:
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o chip8-emulator .

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) .
	./$(BINARY_NAME)

# Run with a specific ROM
run-rom:
	@if [ -z "$(ROM)" ]; then \
		echo "Usage: make run-rom ROM=\"games/filename.ch8\""; \
		echo "Available ROMs:"; \
		ls -1 games/*.ch8 2>/dev/null || echo "No .ch8 files found in games/"; \
		exit 1; \
	fi
	$(GOBUILD) -o $(BINARY_NAME) .
	./$(BINARY_NAME) "$(ROM)"

# Run with race detection
run-race:
	$(GOBUILD) -race -o $(BINARY_NAME) .
	./$(BINARY_NAME)

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
	$(RM) $(BINARY_NAME)
	$(RM) chip8-emulator.exe
	$(RM) coverage.out

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

# List available ROMs
list-roms:
	@echo "Available ROMs in games/ directory:"
	@ls -1 games/*.ch8 2>/dev/null || echo "No .ch8 files found in games/"

# Quick test with a simple ROM (if available)
test-rom:
	@if [ -f "games/Pong [Paul Vervalin, 1990].ch8" ]; then \
		echo "Testing with Pong ROM..."; \
		$(MAKE) run-rom ROM="games/Pong [Paul Vervalin, 1990].ch8"; \
	elif [ -f "games/Tetris [Fran Dachille, 1991].ch8" ]; then \
		echo "Testing with Tetris ROM..."; \
		$(MAKE) run-rom ROM="games/Tetris [Fran Dachille, 1991].ch8"; \
	else \
		echo "No test ROMs found. Use 'make list-roms' to see available ROMs."; \
		echo "Then run: make run-rom ROM=\"games/your-rom.ch8\""; \
	fi

# Help target
help:
	@echo "Available targets:"
	@echo "  build          - Build the application for current platform"
	@echo "  build-windows  - Build for Windows (from WSL)"
	@echo "  build-linux    - Build for Linux"
	@echo "  build-mac      - Build for macOS"
	@echo "  build-linux-arm64 - Build for ARM64 Linux"
	@echo "  run            - Build and run the application"
	@echo "  run-rom        - Build and run with specific ROM (ROM=\"games/filename.ch8\")"
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
	@echo "  list-roms      - List available ROMs in games/ directory"
	@echo "  test-rom       - Quick test with a common ROM"
	@echo "  help           - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make run-rom ROM=games/Pong [Paul Vervalin, 1990].ch8"
	@echo "  make run-rom ROM=games/Tetris [Fran Dachille, 1991].ch8"

.PHONY: all build build-windows build-linux build-mac build-linux-arm64 run run-rom run-race test test-coverage clean deps tidy fmt vet lint install-lint dev list-roms test-rom help 