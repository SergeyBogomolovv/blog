// Package service реализует регистрацию, сессии и заявки.
package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"cleanhome/internal/model"
	"cleanhome/pkg/security"
)

var (
	ErrUnauthorized = errors.New("Войдите, чтобы продолжить.")
	ErrCredentials  = errors.New("Логин или пароль некорректен")
	ErrForbidden    = errors.New("Запрос отклонён. Обновите страницу и повторите попытку.")
	ErrInvalid      = errors.New("Некорректный JSON: проверьте формат, размер и поля запроса.")
	ErrConflict     = errors.New("Этот логин уже занят.")
	ErrNotFound     = errors.New("Страница не найдена.")
)

// Repository контракт хранилища пользователей, заявок и сессий.
type Repository interface {
	// Register сохраняет пользователя и защищает уникальность логина.
	Register(ctx context.Context, user *model.User) error
	// UserByLogin находит пользователя по нормализованному логину.
	UserByLogin(ctx context.Context, login string) (*model.User, error)
	// CreateSession связывает хеш токена с пользователем.
	CreateSession(ctx context.Context, sessionID string, userID int64) error
	// Session проверяет срок сессии и возвращает пользователя.
	Session(ctx context.Context, sessionID string) (*model.User, error)
	// DeleteSession отзывает идентификатор сессии.
	DeleteSession(ctx context.Context, sessionID string) error
	// Requests читает только заявки указанного владельца.
	Requests(ctx context.Context, userID int64) ([]model.Request, error)
	// CreateRequest сохраняет проверенную заявку.
	CreateRequest(ctx context.Context, request *model.Request) error
	// Ping проверяет готовность хранилища.
	Ping(ctx context.Context) error
}

// Service правила регистрации, авторизации и работы с заявками.
type Service struct {
	repo      Repository
	dummyHash string
}

// New создаёт сервис и фиктивный хеш для проверки неизвестных логинов.
func New(repo Repository) (*Service, error) {
	dummyHash, err := security.HashPassword("Dummy-password1!")
	if err != nil {
		return nil, err
	}

	return &Service{
		repo:      repo,
		dummyHash: dummyHash,
	}, nil
}

// Register создаёт пользователя после проверки регистрационных данных.
func (s *Service) Register(ctx context.Context, input Registration) error {
	if fields := input.Validate(); len(fields) > 0 {
		return fields
	}

	passwordHash, err := security.HashPassword(input.Password)
	if err != nil {
		return err
	}

	user := model.User{
		Surname:        input.Surname,
		Name:           input.Name,
		Patronymic:     input.Patronymic,
		Login:          input.Login,
		Phone:          input.Phone,
		Email:          input.Email,
		PasswordHash:   passwordHash,
		ConsentVersion: ConsentVersion,
	}

	return s.repo.Register(ctx, &user)
}

// Login проверяет формат учётных данных и хеш пароля.
func (s *Service) Login(ctx context.Context, login, password string) (*model.User, error) {
	login = strings.TrimSpace(login)
	if !validText(login, loginPattern) || !validPassword(password) {
		return nil, ErrCredentials
	}
	user, err := s.repo.UserByLogin(ctx, login)
	if errors.Is(err, sql.ErrNoRows) {
		// Несуществующий логин также проходит scrypt, чтобы не раскрывать его наличия.
		security.CheckPassword(s.dummyHash, password)
		return nil, ErrCredentials
	}
	if err != nil {
		return nil, err
	}
	if !security.CheckPassword(user.PasswordHash, password) {
		return nil, ErrCredentials
	}

	return user, nil
}

// CreateSession сохраняет идентификатор сессии и её владельца.
func (s *Service) CreateSession(ctx context.Context, sessionID string, userID int64) error {
	return s.repo.CreateSession(ctx, sessionID, userID)
}

// Authorize возвращает пользователя только для действующей сессии.
func (s *Service) Authorize(ctx context.Context, sessionID string) (*model.User, error) {
	user, err := s.repo.Session(ctx, sessionID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUnauthorized
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Logout отзывает сессию в хранилище.
func (s *Service) Logout(ctx context.Context, sessionID string) error {
	return s.repo.DeleteSession(ctx, sessionID)
}

// Requests читает заявки пользователя в порядке от новых к старым.
func (s *Service) Requests(ctx context.Context, userID int64) ([]model.Request, error) {
	return s.repo.Requests(ctx, userID)
}

// CreateRequest проверяет поля заявки и передаёт её в хранилище.
func (s *Service) CreateRequest(
	ctx context.Context,
	userID int64,
	input RequestInput,
) (*model.Request, error) {
	if fields := input.Validate(time.Now()); len(fields) > 0 {
		return nil, fields
	}

	request := model.Request{
		UserID:       userID,
		Service:      input.Service,
		Address:      input.Address,
		CleaningDate: input.Date,
	}
	if err := s.repo.CreateRequest(ctx, &request); err != nil {
		return nil, err
	}

	return &request, nil
}

// Health проверяет готовность хранилища к обработке запросов.
func (s *Service) Health(ctx context.Context) error {
	return s.repo.Ping(ctx)
}
