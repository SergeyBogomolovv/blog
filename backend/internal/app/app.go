// Package app собирает соединение с БД, сервисы и HTTP-сервер.
package app

import (
	"context"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"cleanhome/internal/config"
	"cleanhome/internal/handler"
	"cleanhome/internal/repository"
	"cleanhome/internal/service"
	"cleanhome/pkg/security"
)

// App ресурсы приложения и настроенный HTTP-сервер.
type App struct {
	DB      *sqlx.DB
	Service *service.Service
	Server  *http.Server
}

// OpenDB проверяет соединение и ограничивает размер пула.
func OpenDB(cfg config.Config) (*sqlx.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, "postgres", cfg.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}

// New собирает шифрование, хранилище, бизнес-логику и HTTP-сервер.
func New(cfg config.Config) (*App, error) {
	dataCipher, err := security.NewDataCipher(cfg.DataKey, cfg.IndexKey)
	if err != nil {
		return nil, err
	}
	db, err := OpenDB(cfg)
	if err != nil {
		return nil, err
	}

	repo := repository.New(db, dataCipher)
	svc, err := service.New(repo)
	if err != nil {
		db.Close()
		return nil, err
	}

	server := &http.Server{
		Addr:              cfg.Address,
		Handler:           handler.New(svc, cfg),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    16 * 1024,
	}

	return &App{
		DB:      db,
		Service: svc,
		Server:  server,
	}, nil
}
