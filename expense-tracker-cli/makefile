# Makefile for building and running the expense-tracker-cli project

# Variables
BINARY_NAME=expense-tracker-cli
SOURCE_DIR=.

# Build the binary
build:
	go build -o $(BINARY_NAME) $(SOURCE_DIR)

# Run the binary
run: build
	./$(BINARY_NAME)

# Clean up build artifacts
clean:
	rm -f $(BINARY_NAME)

# Run tests
test:
	go test ./...

# Format the code
fmt:
	go fmt ./...

# Lint the code
lint:
	golangci-lint run

.PHONY: build run clean test fmt lint