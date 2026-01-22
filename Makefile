.PHONY: build run clean test setup docker-build docker-up docker-down help

# Default target
.DEFAULT_GOAL := help

# Variables
BINARY_NAME=signal-bot
DOCKER_IMAGE=signal-bot
DOCKER_COMPOSE=docker-compose

## help: Display this help message
help:
	@echo "Signal Bot - Available commands:"
	@echo ""
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## setup: Run the setup script to configure the bot
setup:
	@./setup.sh

## build: Build the bot binary
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME)
	@echo "Build complete: $(BINARY_NAME)"

## run: Run the bot locally
run: build
	@echo "Starting $(BINARY_NAME)..."
	@./$(BINARY_NAME)

## clean: Remove built binaries and temporary files
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@go clean
	@echo "Clean complete"

## test: Run tests (if any)
test:
	@echo "Running tests..."
	@go test -v ./...

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies updated"

## docker-build: Build the Docker image
docker-build:
	@echo "Building Docker image..."
	@$(DOCKER_COMPOSE) build
	@echo "Docker build complete"

## docker-up: Start the bot using Docker Compose
docker-up:
	@echo "Starting bot with Docker Compose..."
	@$(DOCKER_COMPOSE) up -d
	@echo "Bot is running in the background"

## docker-down: Stop the bot using Docker Compose
docker-down:
	@echo "Stopping bot..."
	@$(DOCKER_COMPOSE) down
	@echo "Bot stopped"

## docker-logs: View bot logs
docker-logs:
	@$(DOCKER_COMPOSE) logs -f

## register: Register Signal number (requires PHONE env var)
register:
ifndef PHONE
	@echo "Error: PHONE variable is required"
	@echo "Usage: make register PHONE=+1234567890"
	@exit 1
endif
	@signal-cli -a $(PHONE) register
	@echo "Check your phone for the verification code"

## verify: Verify Signal number (requires PHONE and CODE env vars)
verify:
ifndef PHONE
	@echo "Error: PHONE variable is required"
	@echo "Usage: make verify PHONE=+1234567890 CODE=123-456"
	@exit 1
endif
ifndef CODE
	@echo "Error: CODE variable is required"
	@echo "Usage: make verify PHONE=+1234567890 CODE=123-456"
	@exit 1
endif
	@signal-cli -a $(PHONE) verify $(CODE)
	@echo "Verification complete!"
