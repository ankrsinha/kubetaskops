.PHONY: help db-up db-down db-logs dev-api dev-web

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

db-up: ## Start local PostgreSQL database container
	docker compose -f backend/config/docker/docker-compose.dev.yml up -d

db-down: ## Stop local PostgreSQL database container
	docker compose -f backend/config/docker/docker-compose.dev.yml down

db-logs: ## View database container logs
	docker compose -f backend/config/docker/docker-compose.dev.yml logs -f postgres

dev-api: ## Run backend Go API server in dev mode
	cd backend && go run cmd/api/main.go

dev-web: ## Run React web frontend in dev mode
	cd web && npm run dev
