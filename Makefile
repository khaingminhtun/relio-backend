APP_NAME := relio-backend

# ============================================================
# Database
# ============================================================

DB_CONTAINER := production-go-postgres

DB_HOST := localhost
DB_PORT := 5433
DB_USER := production
DB_PASSWORD := production
DB_NAME := production_api

TEST_DB_NAME := production_api_test

DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

TEST_DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(TEST_DB_NAME)?sslmode=disable

MIGRATIONS := ./migrations


# ============================================================
# Docker
# ============================================================

.PHONY: docker-up
docker-up:
	docker compose up -d

.PHONY: docker-down
docker-down:
	docker compose down

.PHONY: docker-stop
docker-stop:
	docker compose stop

.PHONY: docker-restart
docker-restart:
	docker compose restart

.PHONY: docker-logs
docker-logs:
	docker compose logs -f

.PHONY: docker-ps
docker-ps:
	docker compose ps


# ============================================================
# Development Database
# ============================================================

.PHONY: db-shell
db-shell:
	docker exec -it $(DB_CONTAINER) psql \
		-U $(DB_USER) \
		-d $(DB_NAME)

.PHONY: db-tables
db-tables:
	docker exec -it $(DB_CONTAINER) psql \
		-U $(DB_USER) \
		-d $(DB_NAME) \
		-c "\dt"

.PHONY: db-drop-all
db-drop-all:
	docker exec -it $(DB_CONTAINER) psql \
		-U $(DB_USER) \
		-d $(DB_NAME) \
		-c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"

.PHONY: db-reset
db-reset:
	docker compose down -v
	docker compose up -d


.PHONY: db-clear-data
db-clear-data:
	docker exec -it $(DB_CONTAINER) psql \
		-U $(DB_USER) \
		-d $(DB_NAME) \
		-c "DO $$$$ DECLARE r RECORD; BEGIN FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP EXECUTE 'TRUNCATE TABLE public.' || quote_ident(r.tablename) || ' RESTART IDENTITY CASCADE'; END LOOP; END $$$$;"
# ============================================================
# Development Database Migrations
# ============================================================

.PHONY: migrate-create
migrate-create:
	goose -dir $(MIGRATIONS) create $(name) sql

.PHONY: migrate-up
migrate-up:
	goose -dir $(MIGRATIONS) postgres "$(DB_URL)" up

.PHONY: migrate-down
migrate-down:
	goose -dir $(MIGRATIONS) postgres "$(DB_URL)" down

.PHONY: migrate-status
migrate-status:
	goose -dir $(MIGRATIONS) postgres "$(DB_URL)" status

.PHONY: migrate-reset
migrate-reset:
	goose -dir $(MIGRATIONS) postgres "$(DB_URL)" reset


# ============================================================
# Test Database
# ============================================================

.PHONY: test-db-create
test-db-create:
	docker exec $(DB_CONTAINER) \
		psql -U $(DB_USER) -d postgres \
		-c "CREATE DATABASE $(TEST_DB_NAME);"

.PHONY: test-db-drop
test-db-drop:
	docker exec $(DB_CONTAINER) \
		psql -U $(DB_USER) -d postgres \
		-c "DROP DATABASE IF EXISTS $(TEST_DB_NAME);"

.PHONY: test-db-shell
test-db-shell:
	docker exec -it $(DB_CONTAINER) \
		psql -U $(DB_USER) -d $(TEST_DB_NAME)

.PHONY: test-db-tables
test-db-tables:
	docker exec $(DB_CONTAINER) \
		psql -U $(DB_USER) -d $(TEST_DB_NAME) \
		-c "\dt"


# ============================================================
# Test Database Migrations
# ============================================================

.PHONY: test-migrate-up
test-migrate-up:
	goose -dir $(MIGRATIONS) postgres "$(TEST_DB_URL)" up

.PHONY: test-migrate-down
test-migrate-down:
	goose -dir $(MIGRATIONS) postgres "$(TEST_DB_URL)" down

.PHONY: test-migrate-status
test-migrate-status:
	goose -dir $(MIGRATIONS) postgres "$(TEST_DB_URL)" status

.PHONY: test-migrate-reset
test-migrate-reset:
	goose -dir $(MIGRATIONS) postgres "$(TEST_DB_URL)" reset


# ============================================================
# Tests
# ============================================================

.PHONY: test
test:
	go test ./...

.PHONY: test-verbose
test-verbose:
	go test ./... -v

.PHONY: test-unit
test-unit:
	go test ./internal/features/... -v

.PHONY: test-integration
test-integration:
	go test ./internal/features/... -v

.PHONY: test-coverage
test-coverage:
	go test ./... -cover

.PHONY: test-race
test-race:
	go test ./... -race


# ============================================================
# Redis
# ============================================================

.PHONY: redis-up
redis-up:
	docker compose up -d redis

.PHONY: redis-down
redis-down:
	docker compose stop redis

.PHONY: redis-restart
redis-restart:
	docker compose restart redis

.PHONY: redis-logs
redis-logs:
	docker compose logs -f redis

.PHONY: redis-status
redis-status:
	docker compose ps redis

.PHONY: redis-cli
redis-cli:
	docker exec -it production-redis redis-cli

.PHONY: redis-ping
redis-ping:
	docker exec production-redis redis-cli PING

.PHONY: redis-flush
redis-flush:
	docker exec production-redis redis-cli FLUSHDB

.PHONY: redis-keys
redis-keys:
	docker exec production-redis redis-cli KEYS '*'

.PHONY: redis-info
redis-info:
	docker exec production-redis redis-cli INFO


# ============================================================
# Application
# ============================================================

.PHONY: run
run:
	go run ./cmd/api

.PHONY: build
build:
	go build -o bin/api ./cmd/api

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...