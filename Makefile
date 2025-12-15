SERVICE := sumni-finance-backend
# Database Config Variables (Use environment variables if set, otherwise use default placeholders)
POSTGRES_USER ?= sumni
POSTGRES_PASSWORD ?= sumni
POSTGRES_HOST ?= localhost
POSTGRES_DATABASE ?= sumni-finance
POSTGRES_PORT ?= 5432

DB_URL := postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable
MIGRATE_PATH := db/migrations

.PHONY: test test-ci dev stop lint migrate-create migrate-up migrate-down

test:
	@./scripts/test.sh .e2e.env

test-ci: 
	@CI=true ./scripts/test.sh .e2e.env

dev:
	DEBUG=$(DEBUG) docker compose up --build $(SERVICE) -d	

logs:
	docker logs -f $(SERVICE)

stop:
	docker compose down $(SERVICE)
	
lint:
	golangci-lint run

# Applies all pending migrations (Forward)
migrate-create:
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(NAME)

migrate-up:
	migrate -database "$(DB_URL)" -path $(MIGRATE_PATH) up

# Rolls back the last applied migration (Backward)
migrate-down:
	migrate -database "$(DB_URL)" -path $(MIGRATE_PATH) down 1

# Shows the current database migration status
migrate-status:
	@echo "Checking migration status..."
	@migrate -database "$(DB_URL)" -path $(MIGRATE_PATH) version