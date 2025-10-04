# ==========================
# Variables
# ==========================
MIGRATIONS_DIR=./migrations/db
DB_URL?=$(shell grep DATABASE_URL .env | cut -d '=' -f2-)
MIGRATE_CMD=migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)"

# ==========================
# Swagger
# ==========================
.PHONY: swagger-gen
swagger-gen:
	swag init --generalInfo src/app.go --output docs

# ==========================
# Migration Commands
# ==========================
.PHONY: migrate-create
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "❌ Please provide a migration name: make migrate-create name=create_table_users"; \
		exit 1; \
	fi
	migrate create -ext sql -dir $(MIGRATIONS_DIR) $(name)

.PHONY: migrate-up
migrate-up:
	$(MIGRATE_CMD) up

.PHONY: migrate-up-1
migrate-up-1:
	$(MIGRATE_CMD) up 1

.PHONY: migrate-down
migrate-down:
	$(MIGRATE_CMD) down

.PHONY: migrate-down-1
migrate-down-1:
	$(MIGRATE_CMD) down 1

.PHONY: migrate-version
migrate-version:
	$(MIGRATE_CMD) version


# ==========================
# Run specific SQL migration file
# ==========================
# Convert DB_URL for psql (hapus prefix "postgresql://")
PSQL_URL=$(shell echo $(DB_URL) | sed 's/^postgresql:\/\///')

# Extract connection parts for psql
DB_USER=$(shell echo $(DB_URL) | sed -E 's#.*//([^:]+):.*#\1#')
DB_PASS=$(shell echo $(DB_URL) | sed -E 's#.*//[^:]+:([^@]+)@.*#\1#')
DB_HOST=$(shell echo $(DB_URL) | sed -E 's#.*@([^:]+):.*#\1#')
DB_PORT=$(shell echo $(DB_URL) | sed -E 's#.*:([0-9]+)/.*#\1#')
DB_NAME=$(shell echo $(DB_URL) | sed -E 's#.*/([^?]+)\?.*#\1#')

.PHONY: sql-run
sql-run:
	@if [ -z "$(id)" ] || [ -z "$(dir)" ]; then \
		echo "❌ Usage: make sql-run id=20250929072208 dir=up"; \
		exit 1; \
	fi; \
	SQL_FILE=$(MIGRATIONS_DIR)/$(id)_*.$(dir).sql; \
	if [ ! -f $$SQL_FILE ]; then \
		echo "❌ File not found: $$SQL_FILE"; \
		exit 1; \
	fi; \
	echo "▶️ Running $$SQL_FILE"; \
	PGPASSWORD=$(DB_PASS) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -f $$SQL_FILE

# ==========================
# Help
# ==========================
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  swagger-gen       Generate Swagger documentation"
	@echo "  migrate-create    Create a new migration (usage: make migrate-create name=create_table_users)"
	@echo "  migrate-up        Apply all migrations"
	@echo "  migrate-up-1      Apply one migration step"
	@echo "  migrate-down      Rollback all migrations"
	@echo "  migrate-down-1    Rollback one migration step"
	@echo "  migrate-version   Show current migration version"
	@echo "  sql-run           Run a specific SQL migration file (usage: make sql-run id=20250929072208 dir=up)"
