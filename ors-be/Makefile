.PHONY: build run clean migrate-up migrate-down migrate-create

APP_NAME = ors-be
DB_URL ?= postgres://postgres:postgres@localhost:5432/ors?sslmode=disable

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

.PHONY: deps
deps:
	go mod tidy
