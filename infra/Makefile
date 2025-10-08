DOCKER_COMPOSE_FILE := docker-compose.yml
API_DIR := ../backend/core-api

DB_CONTAINER := raitometer_db
DB_USER := raitometer_user
DB_NAME := raitometer_db
DB_URL := "postgres://${DB_USER}:raitometer_password@localhost:5432/${DB_NAME}?sslmode=disable"
MIGRATIONS_PATH := $(API_DIR)/internal/migrations


.PHONY: \
help \
run \
run/api \
infra/up \
infra/down \
db/open \
migrate/up \
migrate/down \
migrate/create \
migrate/force \
logs/ai \
logs/db \


help:
	@echo "Available commands:"
	@echo "  make infra/up  --- Run the containers"
	@echo "  make infra/down  --- Stop the containers"
	@echo "  make db/open  --- Open the DB container"
	@echo "  make migrate/create name=<migration_name>  --- Create a new migration file"
	@echo "  make migrate/up  --- Apply all up migrations"
	@echo "  make migrate/down  --- Apply one down migration"
	@echo "  make migrate/force v=<version>  --- Force migration to a specific version"


infra/up:
	@docker compose -f $(DOCKER_COMPOSE_FILE) up --build -d

infra/down:
	@docker compose -f $(DOCKER_COMPOSE_FILE) down -v

db/open:
	@docker exec -it $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME)

migrate/create:
	@if [ -z "$(name)" ]; then \
		echo "Error: 'name' argument is required. Usage: make migrate/create name=your_migration_name"; \
		exit 1; \
	fi
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

migrate/up:
	@migrate -database "$(DB_URL)" -path $(MIGRATIONS_PATH) up

migrate/down:
	@migrate -database "$(DB_URL)" -path $(MIGRATIONS_PATH) down 1

migrate/force:
	@if [ -z "$(v)" ]; then \
		echo "Error: 'v' (version) argument is required. Usage: make migrate/force v=<version_number>"; \
		exit 1; \
	fi
	migrate -database $(DB_URL) -path $(MIGRATIONS_PATH) force $(v)

logs/ai:
	docker logs raitometer_ai

logs/db:
	docker logs raitometer_db