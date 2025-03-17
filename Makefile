.PHONY: run build test lint clean docker-up docker-down

# Variables
APP_NAME=codestacker-api.exe
DOCKER_COMPOSE=deployments/docker-compose.yml

# Run the application locally
run:
	@go run cmd/main.go

# Build the application
build:
	@go build -o $(APP_NAME) ./cmd/main.go

# Run unit tests
test:
	@go test ./... -v

# Lint the code (install golangci-lint first)
lint:
	@golangci-lint run ./...

# Remove compiled binaries
clean:
	@rm -rf $(APP_NAME)

# Start Docker containers
docker-up:
	@docker-compose -f $(DOCKER_COMPOSE) up --build -d

# Stop and remove Docker containers
docker-down:
	@docker-compose -f $(DOCKER_COMPOSE) down
