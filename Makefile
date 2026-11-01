COMPOSE ?= docker compose
.DEFAULT_GOAL := help

.PHONY: help env up down build restart logs ps config db-up db-shell test swagger migrate-up migrate-down dev dev-down

help: ## Show available commands
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z_-]+:.*## / {printf "  %-12s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

env: ## Create .env from env_example if it does not exist
	@test -f .env || cp env_example .env

up: env ## Start the API with Air and PostgreSQL in the background
	@mkdir -p tmp
	$(COMPOSE) up -d --build --wait

dev: up ## Alias for up: start development with automatic rebuilds

dev-down: down ## Alias for down: stop development containers

down: ## Stop and remove containers, preserving database data
	$(COMPOSE) down

build: ## Build the application image
	$(COMPOSE) build app

restart: ## Restart the application
	$(COMPOSE) restart app

logs: ## Follow logs (optionally: make logs SERVICE=app)
	$(COMPOSE) logs -f $(SERVICE)

ps: ## Show service status
	$(COMPOSE) ps

config: ## Validate Compose configuration
	$(COMPOSE) config --quiet

db-up: env ## Start PostgreSQL only and wait for readiness
	$(COMPOSE) up -d --wait postgres

db-shell: ## Open psql in the database container
	$(COMPOSE) exec postgres sh -c 'exec psql -U "$$POSTGRES_USER" -d "$$POSTGRES_DB"'

test: ## Run Go tests in Docker
	docker build --target test --progress=plain .

swagger: ## Generate Swagger docs using the pinned swag CLI (requires Go)
	go run github.com/swaggo/swag/cmd/swag@v1.8.12 init -g main.go -d cmd/api,internal/order,internal/warehouse,internal/product,internal/helpers -o docs

migrate-up: db-up ## Apply pending SQL migrations
	@mkdir -p tmp
	$(COMPOSE) run --rm --no-deps --build --entrypoint /app/migrate app up

migrate-down: db-up ## Revert the latest migration (stop the API first)
	@mkdir -p tmp
	$(COMPOSE) run --rm --no-deps --build --entrypoint /app/migrate app down
