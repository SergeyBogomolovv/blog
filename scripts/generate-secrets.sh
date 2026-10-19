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
POSTGRES_DB=studyblog
POSTGRES_PASSWORD=$(secret)
DB_USER_PASSWORD=$(secret)
DB_ADMIN_PASSWORD=$(secret)
JWT_SECRET=$(secret)
ENCRYPTION_KEY=$(secret)
LOGIN_INDEX_KEY=$(secret)
DB_HOST=127.0.0.1
DB_PORT=5432
SERVER_ADDR=127.0.0.1:8080
APP_ORIGIN=http://localhost:3000
COOKIE_SECURE=false
EOF
echo 'Создан .env с правами 600. Не публикуйте этот файл.'
