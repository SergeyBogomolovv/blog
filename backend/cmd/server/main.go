// Команда server запускает HTTP API и корректно завершает активные запросы.
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

	"studyblog/internal/app"
	"studyblog/internal/config"
)

// main запускает сервер и ожидает сигнал завершения.
func main() {
	if err := run(); err != nil {
		slog.Error("Не удалось запустить сервер", "error", err)
		os.Exit(1)
	}
}

// run управляет жизненным циклом сервера и соединений.
func run() error {
	c, err := config.Load(false)
	if err != nil {
		return err
	}

	a, err := app.New(c)
	if err != nil {
		return err
	}

	defer a.DB.Close()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	result := make(chan error, 1)
	go func() {
		slog.Info("Сервер запущен", "address", c.Address)
		result <- a.Server.ListenAndServe()
	}()
	select {
	case err := <-result:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return a.Server.Shutdown(shutdown)
	}
}
