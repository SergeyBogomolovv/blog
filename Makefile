# Команды разработки; фронтенд и бэкенд запускаются в разных терминалах.
.PHONY: help secrets install db-up db-down migrate admin backend frontend build fmt

help:
	@echo "secrets install db-up db-down migrate admin backend frontend build fmt"

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

admin:
	cd backend && go run ./cmd/admin

backend:
	cd backend && go run ./cmd/server

frontend:
	cd frontend && npm run dev

build:
	cd backend && go build -o bin/server ./cmd/server
	cd backend && go build -o bin/admin ./cmd/admin
	cd backend && go build -o bin/migrate ./cmd/migrate
	cd frontend && npm run build

fmt:
	cd backend && gofmt -w cmd internal pkg migrations
	cd frontend && npm run format
