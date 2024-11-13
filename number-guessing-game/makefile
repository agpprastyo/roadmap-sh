.PHONY: all build run clean

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=number-guessing-game
BINARY_UNIX=$(BINARY_NAME)_unix
MAIN_PATH=./cmd/main.go

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

run: build
	./$(BINARY_NAME)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

test:
	$(GOTEST) -v ./...

deps:
	$(GOGET) -u ./...