.PHONY: help init start stop restart status logs reset shell seed test proto proto-auth proto-game proto-packs proto-lobby

# Default target
.DEFAULT_GOAL := help

# Colors for output
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

help: ## Show this help message
	@echo "$(BLUE)SIGame 2.0 - Infrastructure Management$(NC)"
	@echo ""
	@echo "$(GREEN)Available commands:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}'

init: ## Initialize environment (copy .env.example to .env)
	@echo "$(BLUE)Initializing environment...$(NC)"
	@if [ ! -f .env ]; then \
		cp env.example .env; \
		echo "$(GREEN)✓ Created .env file from env.example$(NC)"; \
		echo "$(YELLOW)! Please review and update .env with your configuration$(NC)"; \
	else \
		echo "$(YELLOW)! .env file already exists, skipping...$(NC)"; \
	fi

start: ## Start all infrastructure services
	@echo "$(BLUE)Starting SIGame infrastructure...$(NC)"
	@docker compose up -d
	@echo "$(GREEN)✓ Infrastructure started successfully$(NC)"
	@echo ""
	@echo "$(YELLOW)Available endpoints:$(NC)"
	@echo "  • PostgreSQL Auth:    localhost:5432"
	@echo "  • PostgreSQL Lobby:   localhost:5433"
	@echo "  • PostgreSQL Packs:   localhost:5434"
	@echo "  • Redis:              localhost:6379"
	@echo "  • Kafka:              localhost:9092"
	@echo "  • MinIO Console:      http://localhost:9001"
	@echo "  • Prometheus:         http://localhost:9090"
	@echo "  • Grafana:            http://localhost:3000 (admin/admin)"

stop: ## Stop all infrastructure services
	@echo "$(BLUE)Stopping SIGame infrastructure...$(NC)"
	@docker compose stop
	@echo "$(GREEN)✓ Infrastructure stopped$(NC)"

down: ## Stop and remove all containers
	@echo "$(BLUE)Stopping and removing containers...$(NC)"
	@docker compose down
	@echo "$(GREEN)✓ Containers removed$(NC)"

restart: ## Restart all services or specific service (use: make restart service=<name>)
	@if [ -z "$(service)" ]; then \
		echo "$(BLUE)Restarting all services...$(NC)"; \
		docker compose restart; \
		echo "$(GREEN)✓ All services restarted$(NC)"; \
	else \
		echo "$(BLUE)Restarting $(service)...$(NC)"; \
		docker compose restart $(service); \
		echo "$(GREEN)✓ $(service) restarted$(NC)"; \
	fi

status: ## Show status of all containers
	@echo "$(BLUE)Container Status:$(NC)"
	@docker compose ps

logs: ## Show logs for specific service (use: make logs service=<name> [follow=true])
	@if [ -z "$(service)" ]; then \
		echo "$(RED)✗ Please specify a service: make logs service=<name>$(NC)"; \
		echo "$(YELLOW)Available services:$(NC)"; \
		echo "  postgres-auth, postgres-lobby, postgres-packs, redis,"; \
		echo "  zookeeper, kafka, minio, prometheus, grafana"; \
		exit 1; \
	fi
	@if [ "$(follow)" = "true" ]; then \
		docker compose logs -f $(service); \
	else \
		docker compose logs --tail=100 $(service); \
	fi

shell: ## Open shell in specific container (use: make shell service=<name>)
	@if [ -z "$(service)" ]; then \
		echo "$(RED)✗ Please specify a service: make shell service=<name>$(NC)"; \
		exit 1; \
	fi
	@echo "$(BLUE)Opening shell in $(service)...$(NC)"
	@docker compose exec $(service) /bin/sh || docker compose exec $(service) /bin/bash

reset: ## DANGER: Stop all services and remove all volumes (data will be lost!)
	@echo "$(RED)⚠ WARNING: This will delete ALL data!$(NC)"
	@echo "$(YELLOW)Press Ctrl+C to cancel, or wait 5 seconds to continue...$(NC)"
	@sleep 5
	@echo "$(BLUE)Stopping and removing all containers and volumes...$(NC)"
	@docker compose down -v
	@echo "$(GREEN)✓ Reset complete - all data has been removed$(NC)"

clean: ## Remove stopped containers and unused volumes
	@echo "$(BLUE)Cleaning up Docker resources...$(NC)"
	@docker compose down --remove-orphans
	@docker volume prune -f
	@echo "$(GREEN)✓ Cleanup complete$(NC)"

build: ## Build/rebuild all services
	@echo "$(BLUE)Building services...$(NC)"
	@docker compose build
	@echo "$(GREEN)✓ Build complete$(NC)"

pull: ## Pull latest images
	@echo "$(BLUE)Pulling latest images...$(NC)"
	@docker compose pull
	@echo "$(GREEN)✓ Images updated$(NC)"

seed: ## Load test data into databases (TODO: implement seed scripts)
	@echo "$(YELLOW)⚠ Seed functionality not yet implemented$(NC)"
	@echo "$(BLUE)This will load test data:$(NC)"
	@echo "  • 5 test users"
	@echo "  • 3 test .siq packs"
	@echo "  • 2 test game rooms"

test: ## Run integration tests (TODO: implement tests)
	@echo "$(YELLOW)⚠ Test functionality not yet implemented$(NC)"

# Proto generation

proto: ## Generate proto files for all services
	@echo "$(BLUE)Generating proto files...$(NC)"
	@bash scripts/generate-proto.sh
	@echo "$(GREEN)✓ Proto generation complete$(NC)"

proto-auth: ## Generate proto files for auth service
	@echo "$(BLUE)Generating proto for auth service...$(NC)"
	@cd services/auth && bash generate-proto.sh

proto-game: ## Generate proto files for game service
	@echo "$(BLUE)Generating proto for game service...$(NC)"
	@cd services/game && bash generate-proto.sh

proto-packs: ## Generate proto files for packs service
	@echo "$(BLUE)Generating proto for packs service...$(NC)"
	@cd services/packs && bash generate-proto.sh

proto-lobby: ## Generate proto files for lobby service (via Gradle)
	@echo "$(BLUE)Generating proto for lobby service...$(NC)"
	@cd services/lobby && ./gradlew generateProto

ps: ## Alias for status
	@$(MAKE) status

# Infrastructure-specific commands

postgres-auth: ## Connect to auth PostgreSQL database
	@docker compose exec postgres-auth psql -U authuser -d auth_db

postgres-lobby: ## Connect to lobby PostgreSQL database
	@docker compose exec postgres-lobby psql -U lobbyuser -d lobby_db

postgres-packs: ## Connect to packs PostgreSQL database
	@docker compose exec postgres-packs psql -U packsuser -d packs_db

redis-cli: ## Connect to Redis CLI
	@docker compose exec redis redis-cli

kafka-topics: ## List Kafka topics
	@docker compose exec kafka kafka-topics --bootstrap-server localhost:9092 --list

kafka-create-topic: ## Create a Kafka topic (use: make kafka-create-topic topic=<name> partitions=3 replication=1)
	@if [ -z "$(topic)" ]; then \
		echo "$(RED)✗ Please specify a topic name: make kafka-create-topic topic=<name>$(NC)"; \
		exit 1; \
	fi
	@docker compose exec kafka kafka-topics --bootstrap-server localhost:9092 \
		--create --topic $(topic) \
		--partitions $(or $(partitions),3) \
		--replication-factor $(or $(replication),1)

# Monitoring

metrics: ## Open Prometheus in browser
	@echo "$(BLUE)Opening Prometheus...$(NC)"
	@open http://localhost:9090 || xdg-open http://localhost:9090 || echo "Please open http://localhost:9090"

dashboard: ## Open Grafana in browser
	@echo "$(BLUE)Opening Grafana...$(NC)"
	@open http://localhost:3000 || xdg-open http://localhost:3000 || echo "Please open http://localhost:3000"

minio: ## Open MinIO console in browser
	@echo "$(BLUE)Opening MinIO console...$(NC)"
	@open http://localhost:9001 || xdg-open http://localhost:9001 || echo "Please open http://localhost:9001"

# Health checks

health: ## Check health of all services
	@echo "$(BLUE)Checking service health...$(NC)"
	@echo ""
	@echo "$(YELLOW)PostgreSQL Auth:$(NC)"
	@docker compose exec -T postgres-auth pg_isready -U authuser || echo "$(RED)✗ Not ready$(NC)"
	@echo ""
	@echo "$(YELLOW)PostgreSQL Lobby:$(NC)"
	@docker compose exec -T postgres-lobby pg_isready -U lobbyuser || echo "$(RED)✗ Not ready$(NC)"
	@echo ""
	@echo "$(YELLOW)PostgreSQL Packs:$(NC)"
	@docker compose exec -T postgres-packs pg_isready -U packsuser || echo "$(RED)✗ Not ready$(NC)"
	@echo ""
	@echo "$(YELLOW)Redis:$(NC)"
	@docker compose exec -T redis redis-cli ping || echo "$(RED)✗ Not ready$(NC)"
	@echo ""
	@echo "$(YELLOW)Kafka:$(NC)"
	@docker compose exec -T kafka kafka-broker-api-versions --bootstrap-server localhost:9092 > /dev/null 2>&1 && echo "$(GREEN)✓ Ready$(NC)" || echo "$(RED)✗ Not ready$(NC)"
	@echo ""
	@echo "$(YELLOW)MinIO:$(NC)"
	@curl -f http://localhost:9000/minio/health/live > /dev/null 2>&1 && echo "$(GREEN)✓ Ready$(NC)" || echo "$(RED)✗ Not ready$(NC)"

