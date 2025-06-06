# Makefile for Mushroom Classifier Go version

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names
BINARY_NAME=mushroom-classifier
TEST_API_BINARY=test-api

# Build directories
BUILD_DIR=build

# All target
all: deps build

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Build the main application
build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v ./main.go

# Build the test API utility
test-api:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(TEST_API_BINARY) -v ./cmd/test-api/main.go

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Run the main application
run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run the API test
test: test-api
	./$(BUILD_DIR)/$(TEST_API_BINARY)

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Install the application
install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

# Cross compilation targets
build-linux:
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 -v ./main.go

build-windows:
	mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe -v ./main.go

build-darwin:
	mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 -v ./main.go
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 -v ./main.go

# Build all platforms
build-all: build-linux build-windows build-darwin

.PHONY: all deps build test-api clean run test fmt lint install build-linux build-windows build-darwin build-all