// Package repository выполняет параметрические запросы PostgreSQL через Squirrel и sqlx.
package repository

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"studyblog/internal/model"
)

var builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Repository struct {
	DB *sqlx.DB
}

// New связывает репозиторий с пулом соединений.
func New(db *sqlx.DB) *Repository {
	return &Repository{DB: db}
}

// account выбирает зашифрованную учётную запись по параметрическому условию.
func (r *Repository) account(ctx context.Context, where sq.Eq) (*model.Account, error) {
	q, args, err := builder.
		Select("*").
		From("users").
		Where(where).
		ToSql()
	if err != nil {
		return nil, err
	}

	a := &model.Account{}
	err = r.DB.GetContext(ctx, a, q, args...)
	return a, err
}

// Account загружает учётную запись по идентификатору.
func (r *Repository) Account(ctx context.Context, id int64) (*model.Account, error) {
	return r.account(ctx, sq.Eq{"id": id})
}

// ByLogin ищет логин по индексу HMAC.
func (r *Repository) ByLogin(ctx context.Context, index string) (*model.Account, error) {
	return r.account(ctx, sq.Eq{"login_index": index})
}

// CreateAccount атомарно создаёт аккаунт; поле role доступно только сервисной роли admin.
func (r *Repository) CreateAccount(ctx context.Context, a *model.Account, admin bool) error {
	columns := []string{"login_index", "login_enc", "email_enc", "name_enc", "bio_enc", "password_hash"}
	values := []any{a.LoginIndex, a.LoginEnc, a.EmailEnc, a.NameEnc, a.BioEnc, a.PasswordHash}
	if admin {
		columns = append(columns, "role", "consent_at")
		values = append(values, "admin", nil)
	}

	q, args, err := builder.
		Insert("users").
		Columns(columns...).
		Values(values...).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}

	return r.DB.GetContext(ctx, &a.ID, q, args...)
}

// execute подготавливает и выполняет запрос с привязанными параметрами.
func (r *Repository) execute(ctx context.Context, b sq.Sqlizer) (int64, error) {
	q, args, err := b.ToSql()
	if err != nil {
		return 0, err
	}

	stmt, err := r.DB.PreparexContext(ctx, q)
	if err != nil {
		return 0, err
	}

	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

// Profile обновляет зашифрованные поля публичного профиля.
func (r *Repository) Profile(ctx context.Context, id int64, name, bio string) error {
	_, err := r.execute(ctx, builder.
		Update("users").
		Set("name_enc", name).
		Set("bio_enc", bio).
		Where(sq.Eq{"id": id}),
	)
	return err
}

// Revoke отзывает все выданные токены учётной записи.
func (r *Repository) Revoke(ctx context.Context, id int64) error {
	_, err := r.execute(ctx, builder.
		Update("users").
		Set("token_version", sq.Expr("token_version + 1")).
		Where(sq.Eq{"id": id}),
	)
	return err
}

// Ban меняет блокировку обычного пользователя и отзывает его токены.
func (r *Repository) Ban(ctx context.Context, id int64, banned bool) error {
	n, err := r.execute(ctx, builder.
		Update("users").
		Set("banned", banned).
		Set("token_version", sq.Expr("token_version + 1")).
		Where(sq.Eq{"id": id, "role": "user"}),
	)
	if err == nil && n == 0 {
		return sql.ErrNoRows
	}

	return err
}

// Users возвращает учётные записи для страницы модерации.
func (r *Repository) Users(ctx context.Context) ([]model.Account, error) {
	query, args, err := builder.
		Select("*").
		From("users").
		OrderBy("id DESC").
		Limit(200).
		ToSql()
	if err != nil {
		return nil, err
	}

	out := []model.Account{}
	err = r.DB.SelectContext(ctx, &out, query, args...)
	return out, err
}

// postSelect собирает общую выборку постов с коррелированным подсчётом комментариев.
func postSelect() sq.SelectBuilder {
	commentCount := sq.Select("count(*)").
		From("comments").
		Where(sq.Expr("comments.post_id = posts.id"))

	return builder.
		Select("posts.*").
		Column(sq.Alias(commentCount, "comment_count")).
		From("posts")
}

// Posts возвращает страницу публикаций с необязательным фильтром по автору.
func (r *Repository) Posts(ctx context.Context, author int64, offset uint64) ([]model.Post, error) {
	b := postSelect().
		OrderBy("id DESC").
		Limit(20).
		Offset(offset)
	if author > 0 {
		b = b.Where(sq.Eq{"author_id": author})
	}

	q, args, err := b.ToSql()
	if err != nil {
		return nil, err
	}

	out := []model.Post{}
	err = r.DB.SelectContext(ctx, &out, q, args...)
	return out, err
}

// Post получает публикацию и число комментариев.
func (r *Repository) Post(ctx context.Context, id int64) (*model.Post, error) {
	query, args, err := postSelect().
		Where(sq.Eq{"posts.id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	p := &model.Post{}
	err = r.DB.GetContext(ctx, p, query, args...)
	return p, err
}

// CreatePost сохраняет проверенную публикацию.
func (r *Repository) CreatePost(ctx context.Context, author int64, title, body string) (int64, error) {
	q, args, err := builder.
		Insert("posts").
		Columns("author_id", "title", "body").
		Values(author, title, body).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, err
	}

	var id int64
	err = r.DB.GetContext(ctx, &id, q, args...)
	return id, err
}

// DeletePost удаляет публикацию после проверки прав в сервисе.
func (r *Repository) DeletePost(ctx context.Context, id int64) error {
	_, err := r.execute(ctx, builder.
		Delete("posts").
		Where(sq.Eq{"id": id}),
	)
	return err
}

// Comments возвращает комментарии в порядке публикации.
func (r *Repository) Comments(ctx context.Context, id int64) ([]model.Comment, error) {
	query, args, err := builder.
		Select("*").
		From("comments").
		Where(sq.Eq{"post_id": id}).
		OrderBy("id").
		Limit(500).
		ToSql()
	if err != nil {
		return nil, err
	}

	out := []model.Comment{}
	err = r.DB.SelectContext(ctx, &out, query, args...)
	return out, err
}

// CreateComment добавляет проверенный комментарий к публикации.
func (r *Repository) CreateComment(ctx context.Context, post, author int64, body string) error {
	_, err := r.execute(ctx, builder.
		Insert("comments").
		Columns("post_id", "author_id", "body").
		Values(post, author, body),
	)
	return err
}

// Ping проверяет доступность PostgreSQL.
func (r *Repository) Ping(ctx context.Context) error {
	return r.DB.PingContext(ctx)
}
