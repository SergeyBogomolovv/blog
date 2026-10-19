#!/usr/bin/env sh
# Создаёт сервисные роли; структуру таблиц создают миграции Goose.
set -eu
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" \
  -v db_name="$POSTGRES_DB" -v user_password="$DB_USER_PASSWORD" \
  -v admin_password="$DB_ADMIN_PASSWORD" -f /docker-entrypoint-initdb.d/roles.sql.inc
