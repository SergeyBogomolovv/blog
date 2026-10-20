// Команда server запускает API и завершает активные запросы при остановке.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pressly/goose/v3"

	"cleanhome/internal/app"
	"cleanhome/internal/config"
	"cleanhome/migrations"
)

// main запускает команду и сообщает об ошибке завершения.
func main() {
	if err := run(); err != nil {
		slog.Error("Не удалось запустить сервер", "error", err)
		os.Exit(1)
	}
}

// run собирает зависимости и выполняет команду приложения.
func run() error {
	if err := applyMigrations(); err != nil {
		return err
	}

	cfg, err := config.Load(false)
	if err != nil {
		return err
	}

	application, err := app.New(cfg)
	if err != nil {
		return err
	}
	defer application.DB.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	return serve(ctx, application.Server)
}

// applyMigrations применяет встроенные миграции через административную роль.
func applyMigrations() error {
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

	return goose.Up(db.DB, ".")
}

// serve обслуживает HTTP-запросы и корректно завершает работу по сигналу.
func serve(ctx context.Context, server *http.Server) error {
	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("Сервер запущен", "address", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		return server.Shutdown(shutdownCtx)
	}
}
