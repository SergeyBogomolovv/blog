# Фронтенд и Go API работают отдельно; миграции применяются при запуске API.
.PHONY: help secrets install db-up db-down migrate backend frontend build test fmt
help:
	@echo "secrets install db-up db-down migrate backend frontend build test fmt"
secrets:
	sh scripts/generate-secrets.sh
install:
	cd backend && go mod download
	cd frontend && npm ci
db-up:
	docker compose up -d --wait db
db-down:
	docker compose stop db
migrate:
	cd backend && go run ./cmd/migrate up
backend:
	cd backend && go run ./cmd/server
frontend:
	cd frontend && npm run dev
build:
	cd backend && go build -o bin/server ./cmd/server
	cd backend && go build -o bin/migrate ./cmd/migrate
	cd frontend && npm run build
test:
	cd backend && go test ./...
	cd frontend && npm test
fmt:
	cd backend && gofmt -w cmd internal pkg migrations
	cd frontend && npm run format
