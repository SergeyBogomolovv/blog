// Команда migrate применяет и откатывает встроенные миграции Goose.
package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/pressly/goose/v3"

	"cleanhome/internal/app"
	"cleanhome/internal/config"
	"cleanhome/migrations"
)

// main завершает процесс с ошибкой при неудачной миграции.
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка миграции:", err)
		os.Exit(1)
	}
}

// run выполняет up, down или status через отдельную сервисную роль admin.
func run() error {
	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	if len(os.Args) > 2 || (command != "up" && command != "down" && command != "status") {
		return errors.New("использование: migrate [up|down|status]")
	}

	cfg, err := config.Load(true)
	if err != nil {
		return err
	}

	db, err := app.OpenDB(cfg)
	if err != nil {
		return err
	}

	defer db.Close()
	goose.SetBaseFS(migrations.Files)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	return goose.RunContext(context.Background(), command, db.DB, ".")
}
