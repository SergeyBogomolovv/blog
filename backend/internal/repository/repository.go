// Package repository хранит данные в PostgreSQL через Squirrel и sqlx.
package repository

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"cleanhome/internal/model"
	"cleanhome/internal/service"
	"cleanhome/pkg/security"
)

// userColumns исключает служебный индекс логина из выдаваемой модели пользователя.
var userColumns = []string{
	"id", "surname", "name", "patronymic", "login", "phone", "email",
	"password_hash", "consent_at", "consent_version", "created_at",
}

// Repository шифрует персональные данные на границе доступа к PostgreSQL.
type Repository struct {
	db     *sqlx.DB
	sql    sq.StatementBuilderType
	cipher *security.DataCipher
}

// New связывает хранилище с защитой персональных данных.
func New(db *sqlx.DB, dataCipher *security.DataCipher) *Repository {
	return &Repository{
		db:     db,
		sql:    sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		cipher: dataCipher,
	}
}

// Ping проверяет доступность соединения с базой данных.
func (r *Repository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// Register шифрует персональные поля и сохраняет пользователя с индексом логина.
func (r *Repository) Register(ctx context.Context, user *model.User) error {
	encrypted := *user
	if err := r.transformUser(&encrypted, r.cipher.Encrypt); err != nil {
		return err
	}
	query, args, err := r.sql.
		Insert("users").
		Columns(
			"surname",
			"name",
			"patronymic",
			"login",
			"login_index",
			"phone",
			"email",
			"password_hash",
			"consent_version",
		).
		Values(
			encrypted.Surname,
			encrypted.Name,
			encrypted.Patronymic,
			encrypted.Login,
			r.cipher.LoginIndex(user.Login),
			encrypted.Phone,
			encrypted.Email,
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

// UserByLogin находит пользователя по закрытому индексу и расшифровывает данные.
func (r *Repository) UserByLogin(ctx context.Context, login string) (*model.User, error) {
	query, args, err := r.sql.
		Select(userColumns...).
		From("users").
		Where(sq.Eq{"login_index": r.cipher.LoginIndex(login)}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := r.db.GetContext(ctx, &user, query, args...); err != nil {
		return nil, err
	}

	if err := r.transformUser(&user, r.cipher.Decrypt); err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateSession сохраняет идентификатор сессии и её владельца.
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

// Session атомарно продлевает действующую сессию и читает её пользователя.
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
		Select(userColumns...).
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

	if err := r.transformUser(&user, r.cipher.Decrypt); err != nil {
		return nil, err
	}
	return &user, nil
}

// DeleteSession удаляет сессию, прекращая доступ по её токену.
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

// Requests читает заявки пользователя в порядке от новых к старым.
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
	for i := range requests {
		address, err := r.cipher.Decrypt(requests[i].Address, addressContext(userID))
		if err != nil {
			return nil, err
		}
		requests[i].Address = address
	}

	return requests, nil
}

// CreateRequest шифрует адрес и сохраняет заявку текущего пользователя.
func (r *Repository) CreateRequest(ctx context.Context, request *model.Request) error {
	address, err := r.cipher.Encrypt(request.Address, addressContext(request.UserID))
	if err != nil {
		return err
	}
	query, args, err := r.sql.
		Insert("requests").
		Columns("user_id", "service", "address", "cleaning_date").
		Values(request.UserID, request.Service, address, request.CleaningDate).
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

// transformUser преобразует только персональные поля, сохраняя хеш пароля отдельно.
func (r *Repository) transformUser(user *model.User, transform func(string, string) (string, error)) error {
	fields := []struct {
		name  string
		value *string
	}{
		{"surname", &user.Surname},
		{"name", &user.Name},
		{"patronymic", &user.Patronymic},
		{"login", &user.Login},
		{"phone", &user.Phone},
		{"email", &user.Email},
	}
	for _, field := range fields {
		value, err := transform(*field.value, "users."+field.name)
		if err != nil {
			return err
		}
		*field.value = value
	}
	return nil
}

// addressContext связывает шифротекст адреса с владельцем заявки.
func addressContext(userID int64) string {
	return fmt.Sprintf("requests.address:%d", userID)
}
