// Package app собирает соединение с БД, сервисы и HTTP-сервер.
package app

import (
	"context"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"studyblog/internal/config"
	"studyblog/internal/handler"
	"studyblog/internal/repository"
	"studyblog/internal/service"
	"studyblog/pkg/security"
)

type App struct {
	DB      *sqlx.DB
	Service *service.Service
	Server  *http.Server
}

// OpenDB проверяет соединение и ограничивает размер пула.
func OpenDB(c config.Config) (*sqlx.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := sqlx.ConnectContext(ctx, "postgres", c.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	return db, nil
}

// New создаёт приложение с зависимостями для сервера или административной утилиты.
func New(c config.Config) (*App, error) {
	db, err := OpenDB(c)
	if err != nil {
		return nil, err
	}

	v, err := security.NewVault(c.EncryptionKey, c.IndexKey)
	if err != nil {
		db.Close()
		return nil, err
	}

	s, err := service.New(repository.New(db), v)
	if err != nil {
		db.Close()
		return nil, err
	}

	return &App{
		DB:      db,
		Service: s,
		Server: &http.Server{
			Addr:              c.Address,
			Handler:           handler.New(s, c),
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       60 * time.Second,
			MaxHeaderBytes:    16384,
		},
	}, nil
}
