.PHONY: run build test docker-up docker-down

APP_NAME=codestacker-api.exe
DOCKER_COMPOSE=deployments/docker-compose.yml

run:
	@go run cmd/main.go

build:
	@go build -o $(APP_NAME) ./cmd/main.go

test:
	@go test ./... -v

docker-up:
	@docker-compose -f $(DOCKER_COMPOSE) up --build -d

docker-down:
	@docker-compose -f $(DOCKER_COMPOSE) down
