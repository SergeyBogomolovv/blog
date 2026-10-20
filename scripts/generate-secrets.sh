#!/usr/bin/env sh
# Создаёт локальные секреты без перезаписи существующей конфигурации.
set -eu
cd "$(dirname "$0")/.."
if [ -e .env ]; then
  echo '.env уже существует; ключи и пароли базы сохранены.'
  exit 0
fi
umask 077
# Генерирует криптографически случайный секрет в шестнадцатеричном виде.
secret() { openssl rand -hex 32; }
set -C
cat > .env <<EOF
POSTGRES_DB=cleanhome
POSTGRES_PASSWORD=$(secret)
DB_USER_PASSWORD=$(secret)
DB_ADMIN_PASSWORD=$(secret)
SESSION_SECRET=$(secret)
DB_HOST=127.0.0.1
DB_PORT=5432
SERVER_ADDR=127.0.0.1:8080
APP_ORIGIN=http://127.0.0.1:3107
EOF
echo 'Создан .env с правами 600. Не публикуйте этот файл.'
