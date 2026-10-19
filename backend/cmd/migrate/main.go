// Команда migrate применяет и откатывает встроенные миграции Goose.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pressly/goose/v3"

	"studyblog/internal/app"
	"studyblog/internal/config"
	"studyblog/migrations"
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
		return fmt.Errorf("использование: migrate [up|down|status]")
	}

	c, err := config.Load(true)
	if err != nil {
		return err
	}

	db, err := app.OpenDB(c)
	if err != nil {
		return err
	}

	defer db.Close()
	goose.SetBaseFS(migrations.Files)
	if err = goose.SetDialect("postgres"); err != nil {
		return err
	}

	return goose.RunContext(context.Background(), command, db.DB, ".")
}
