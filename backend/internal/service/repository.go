// Контракт хранилища отделяет бизнес-логику от реализации PostgreSQL.
package service

import (
	"context"
	"database/sql"
	"errors"

	"studyblog/internal/model"
)

type Repository interface {
	Account(context.Context, int64) (*model.Account, error)
	ByLogin(context.Context, string) (*model.Account, error)
	CreateAccount(context.Context, *model.Account, bool) error
	Profile(context.Context, int64, string, string) error
	Revoke(context.Context, int64) error
	Ban(context.Context, int64, bool) error
	Users(context.Context) ([]model.Account, error)
	Posts(context.Context, int64, uint64) ([]model.Post, error)
	Post(context.Context, int64) (*model.Post, error)
	CreatePost(context.Context, int64, string, string) (int64, error)
	DeletePost(context.Context, int64) error
	Comments(context.Context, int64) ([]model.Comment, error)
	CreateComment(context.Context, int64, int64, string) error
	Ping(context.Context) error
}

// Authorize сверяет действующую версию сессии и состояние учётной записи.
func (s *Service) Authorize(ctx context.Context, id int64, version int) (*model.Account, error) {
	a, err := s.repo.Account(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUnauthorized
	}

	if err != nil {
		return nil, err
	}

	if a.Banned || a.TokenVersion != version {
		return nil, ErrUnauthorized
	}

	return a, nil
}

// Logout отзывает все сессии пользователя.
func (s *Service) Logout(ctx context.Context, id int64) error {
	return s.repo.Revoke(ctx, id)
}

// Ban проверяет роль инициатора и запрещает блокировку администраторов.
func (s *Service) Ban(ctx context.Context, actor *model.Account, id int64, banned bool) error {
	if actor.Role != "admin" {
		return ErrForbidden
	}

	return s.repo.Ban(ctx, id, banned)
}

// Health проверяет готовность хранилища обслуживать запросы.
func (s *Service) Health(ctx context.Context) error {
	return s.repo.Ping(ctx)
}
