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
)

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

func New(cfg config.Config) (*App, error) {
	db, err := OpenDB(cfg)
	if err != nil {
		return nil, err
	}

	repo := repository.New(db)
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
