COMPOSE=docker compose

.PHONY: help dev dev-build down restart logs logs-app logs-db ps build prod clean reset-db generate-swagger

help: ## Show available commands
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## ' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "%-15s %s\n", $$1, $$2}'

## =========================
## Development
## =========================

dev: ## Run project in development mode (hot reload)
	$(COMPOSE) up -d

dev-build: ## Rebuild dev containers
	$(COMPOSE) up --build -d

down: ## Stop containers
	$(COMPOSE) down

restart: ## Restart containers
	$(COMPOSE) restart

logs: ## Show all logs
	$(COMPOSE) logs -f

logs-app: ## Show app logs
	$(COMPOSE) logs -f app

logs-db: ## Show postgres logs
	$(COMPOSE) logs -f postgres

ps: ## Show running containers
	$(COMPOSE) ps


## =========================
## Production
## =========================

build: ## Build production images
	$(COMPOSE) build

prod: ## Run production containers
	$(COMPOSE) -f docker-compose.yml up --build -d


## =========================
## Utilities
## =========================

clean: ## Remove containers, networks, volumes
	$(COMPOSE) down -v

reset-db: ## Reset database volume
	$(COMPOSE) down -v
	$(COMPOSE) up -d postgres

generate-swagger: ## Generate swagger docs
	cd app && swag init -g ./main.go -o ./docs