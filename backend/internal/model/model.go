// Package model описывает данные клинингового сервиса.
package model

import "time"

type User struct {
	ID             int64     `db:"id" json:"id"`
	Surname        string    `db:"surname" json:"surname"`
	Name           string    `db:"name" json:"name"`
	Patronymic     string    `db:"patronymic" json:"patronymic"`
	Login          string    `db:"login" json:"login"`
	Phone          string    `db:"phone" json:"phone"`
	Email          string    `db:"email" json:"email"`
	PasswordHash   string    `db:"password_hash" json:"-"`
	ConsentAt      time.Time `db:"consent_at" json:"consent_at"`
	ConsentVersion string    `db:"consent_version" json:"consent_version"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

type Request struct {
	ID           int64     `db:"id" json:"id"`
	UserID       int64     `db:"user_id" json:"-"`
	Service      string    `db:"service" json:"service"`
	Address      string    `db:"address" json:"address"`
	CleaningDate string    `db:"cleaning_date" json:"cleaning_date"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type CleaningService struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
