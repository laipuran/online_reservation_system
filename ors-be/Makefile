COMPOSE_FILE = ../docker-compose.yml

.PHONY: build run clean migrate-up migrate-down migrate-create db-up db-down seed setup

APP_NAME = ors-be
DB_URL ?= postgres://postgres:postgres@localhost:5432/ors?sslmode=disable

db-up:
	docker compose up -d db

db-down:
	docker compose down

build:
	go build -o bin/$(APP_NAME) ./cmd/server

run: build
	./bin/$(APP_NAME)

clean:
	rm -rf bin/

migrate-up:
	@echo "Applying migrations..."
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	@echo "Rolling back migrations..."
	migrate -path migrations -database "$(DB_URL)" down 1

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

setup: db-up
	@echo "Waiting for DB to be ready..."
	@sleep 3
	$(MAKE) migrate-up
	$(MAKE) seed

seed:
	@echo "Seeding reference data..."
	docker compose exec -T db psql -U postgres -d ors < db/seed.sql

.PHONY: deps
deps:
	go mod tidy
