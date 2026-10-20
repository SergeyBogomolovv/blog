# Фронтенд и Go API работают отдельно; миграции применяются при запуске API.
.PHONY: help secrets install db-up db-down migrate backend frontend build test fmt
# Выводит доступные команды проекта.
help:
	@echo "secrets install db-up db-down migrate backend frontend build test fmt"
# Создаёт локальную конфигурацию с независимыми случайными ключами.
secrets:
	sh scripts/generate-secrets.sh
# Устанавливает зависимости из зафиксированных версий.
install:
	cd backend && go mod download
	cd frontend && npm ci
# Запускает PostgreSQL и ожидает готовности.
db-up:
	docker compose up -d --wait db
# Останавливает PostgreSQL с сохранением данных.
db-down:
	docker compose stop db
# Применяет SQL-миграции административной ролью.
migrate:
	cd backend && go run ./cmd/migrate up
# Запускает сервер API.
backend:
	cd backend && go run ./cmd/server
# Запускает интерфейс в режиме разработки.
frontend:
	cd frontend && npm run dev
# Собирает исполняемые файлы Go и интерфейс для эксплуатации.
build:
	cd backend && go build -o bin/server ./cmd/server
	cd backend && go build -o bin/migrate ./cmd/migrate
	cd frontend && npm run build
# Проверяет серверные механизмы и клиентские правила.
test:
	cd backend && go test ./...
	cd frontend && npm test
# Приводит исходники и документацию к единому оформлению.
fmt:
	cd backend && gofmt -w cmd internal pkg migrations
	cd frontend && npm run format
	cd frontend && npx prettier --write ../README.md ../docker-compose.yml ../.prettierrc.json
	cd frontend && npx prettier --write --parser html public/clean-room.svg app/icon.svg
