# ==============================
# Config (overridable)
# ==============================

MIGRATE          ?= migrate
MIGRATIONS_DIR   ?= db/migrations
DB_DRIVER        ?= postgres

DB_USER          ?= admin
DB_PASSWORD      ?= password
DB_HOST          ?= localhost
DB_PORT          ?= 5435
DB_NAME          ?= chat
DB_SSLMODE       ?= disable

DATABASE_URL := $(DB_DRIVER)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# ==============================
# Meta
# ==============================

.PHONY: help migrations-up migrations-down migration-create db-url build rebuild

help:
	@echo ""
	@echo "Available commands:"
	@echo ""
	@echo "  make migration-create name=<name>   Create a new migration"
	@echo "  make migrations-up                  Apply all migrations"
	@echo "  make migrations-down                Roll back last migration"
	@echo "  make db-url                         Print database URL"
	@echo "  make build                          Docker compose up --build"
	@echo "  make rebuild                        Full rebuild (down + build)"
	@echo ""

# ==============================
# Migrations
# ==============================

migration-create:
	@if [ -z "$(name)" ]; then \
		echo "❌ name is required. Example:"; \
		echo "   make migration-create name=create_users_table"; \
		exit 1; \
	fi
	$(MIGRATE) create \
		-ext sql \
		-dir $(MIGRATIONS_DIR) \
		-seq $(name)

migrations-up:
	$(MIGRATE) -database "$(DATABASE_URL)" -path $(MIGRATIONS_DIR) up

migrations-down:
	$(MIGRATE) -database "$(DATABASE_URL)" -path $(MIGRATIONS_DIR) down 1

# ==============================
# Docker
# ==============================

build:
	docker compose up --build

rebuild:
	docker compose down -v
	docker compose up --build

# ==============================
# Debug
# ==============================

db-url:
	@echo $(DATABASE_URL)
