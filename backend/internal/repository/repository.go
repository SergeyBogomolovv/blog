// Package repository хранит данные в PostgreSQL через Squirrel и sqlx.
package repository

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"cleanhome/internal/model"
	"cleanhome/internal/service"
)

type Repository struct {
	db  *sqlx.DB
	sql sq.StatementBuilderType
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		db:  db,
		sql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *Repository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *Repository) Register(ctx context.Context, user *model.User) error {
	query, args, err := r.sql.
		Insert("users").
		Columns(
			"surname",
			"name",
			"patronymic",
			"login",
			"phone",
			"email",
			"password_hash",
			"consent_version",
		).
		Values(
			user.Surname,
			user.Name,
			user.Patronymic,
			user.Login,
			user.Phone,
			user.Email,
			user.PasswordHash,
			user.ConsentVersion,
		).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)

	var postgresError *pq.Error
	if errors.As(err, &postgresError) && postgresError.Code == "23505" {
		return service.ErrConflict
	}

	return err
}

func (r *Repository) UserByLogin(ctx context.Context, login string) (*model.User, error) {
	query, args, err := r.sql.
		Select("*").
		From("users").
		Where("lower(login) = lower(?)", login).
		ToSql()
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := r.db.GetContext(ctx, &user, query, args...); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) CreateSession(ctx context.Context, sessionID string, userID int64) error {
	query, args, err := r.sql.
		Delete("sessions").
		Where("expires_at <= now()").
		ToSql()
	if err != nil {
		return err
	}
	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	// sessionID уже содержит SHA-256, а не значение cookie пользователя.
	query, args, err = r.sql.
		Insert("sessions").
		Columns("id", "user_id").
		Values(sessionID, userID).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *Repository) Session(ctx context.Context, sessionID string) (*model.User, error) {
	// Проверка срока и продление атомарны: истёкшая сессия не возобновляется.
	query, args, err := r.sql.
		Update("sessions").
		Set("expires_at", sq.Expr("now() + interval '30 minutes'")).
		Where(sq.Eq{"id": sessionID}).
		Where("expires_at > now()").
		Suffix("RETURNING user_id").
		ToSql()
	if err != nil {
		return nil, err
	}

	var userID int64
	if err := r.db.GetContext(ctx, &userID, query, args...); err != nil {
		return nil, err
	}

	query, args, err = r.sql.
		Select("*").
		From("users").
		Where(sq.Eq{"id": userID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := r.db.GetContext(ctx, &user, query, args...); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) DeleteSession(ctx context.Context, sessionID string) error {
	query, args, err := r.sql.
		Delete("sessions").
		Where(sq.Eq{"id": sessionID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *Repository) Requests(ctx context.Context, userID int64) ([]model.Request, error) {
	query, args, err := r.sql.
		Select(
			"id",
			"user_id",
			"service",
			"address",
			"to_char(cleaning_date, 'YYYY-MM-DD') AS cleaning_date",
			"created_at",
		).
		From("requests").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("id DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	requests := make([]model.Request, 0)
	if err := r.db.SelectContext(ctx, &requests, query, args...); err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *Repository) CreateRequest(ctx context.Context, request *model.Request) error {
	query, args, err := r.sql.
		Insert("requests").
		Columns("user_id", "service", "address", "cleaning_date").
		Values(request.UserID, request.Service, request.Address, request.CleaningDate).
		Suffix("RETURNING id, created_at").
		ToSql()
	if err != nil {
		return err
	}

	return r.db.QueryRowxContext(ctx, query, args...).Scan(
		&request.ID,
		&request.CreatedAt,
	)
}
