.PHONY: build run clean test fmt vet lint deps tidy help

# Binary name
BINARY_NAME=gitport
BINARY_PATH=bin/$(BINARY_NAME)

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOVET=$(GOCMD) vet

# Build flags
LDFLAGS=-ldflags "-s -w" # Strip debug info and symbol table for smaller binary
BUILD_FLAGS=-v # Verbose output

# Default target
.DEFAULT_GOAL := help

## build: Build the application binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/server
	@echo "Binary built at $(BINARY_PATH)"

## run: Build and run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_PATH)

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf bin/
	@echo "Clean complete"

## test: Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

## fmt: Format Go code
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .
	@echo "Formatting complete"

## vet: Run go vet (static analysis)
vet:
	@echo "Running go vet..."
	$(GOVET) ./...
	@echo "Vet complete"

## lint: Run golangci-lint (if installed)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download

## tidy: Clean up dependencies
tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

## install: Install the binary to $GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GOCMD) install ./cmd/server
	@echo "Installed to $$(go env GOPATH)/bin/$(BINARY_NAME)"

## build-linux: Build for Linux (cross-compilation)
build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o bin/$(BINARY_NAME)-linux ./cmd/server

## build-darwin: Build for macOS (cross-compilation)
build-darwin:
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin ./cmd/server

## build-windows: Build for Windows (cross-compilation)
build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o bin/$(BINARY_NAME)-windows.exe ./cmd/server

## build-all: Build for all platforms
build-all: build-linux build-darwin build-windows
	@echo "Builds complete for all platforms"

## help: Show this help message
help:
	@echo "Available targets:"
	@grep -E '^##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
