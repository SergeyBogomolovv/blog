// Package migrations встраивает SQL-миграции в исполняемый файл.
package migrations

import "embed"

// Files содержит версионированные изменения схемы PostgreSQL.
//
//go:embed *.sql
var Files embed.FS
