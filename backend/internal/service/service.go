// Package service содержит правила регистрации, авторизации и публикации.
package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"net/mail"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"studyblog/internal/model"
	"studyblog/pkg/security"
)

var (
	ErrPasswordShort = errors.New("Пароль должен содержать не менее 10 символов")
	ErrPasswordLong  = errors.New("Пароль слишком длинный. Попробуйте более короткий")
	ErrInvalid       = errors.New("Проверьте формат и заполнение всех полей")
	ErrCredentials   = errors.New("Логин или пароль некорректен")
	ErrConflict      = errors.New("Этот логин уже занят")
	ErrForbidden     = errors.New("Недостаточно прав для этого действия")
	ErrUnauthorized  = errors.New("Войдите в аккаунт заново")
	ErrNotFound      = errors.New("Запись не найдена")
	loginPattern     = regexp.MustCompile(`^[a-z0-9_]{3,24}$`)
)

type Service struct {
	repo      Repository
	Vault     *security.Vault
	dummyHash []byte
}

type Registration struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Consent  bool   `json:"consent"`
}

// New создаёт сервис и запасной хеш для одинаковой проверки неизвестных логинов.
func New(repo Repository, vault *security.Vault) (*Service, error) {
	dummy := make([]byte, 32)
	if _, err := rand.Read(dummy); err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword(dummy, 12)
	return &Service{repo, vault, hash}, err
}

// normalizeText приводит переносы строк из браузерных форм к единому формату.
func normalizeText(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "\r\n", "\n"), "\r", "\n")
}

// validText проверяет длину непустого текста и отсутствие управляющих символов.
func validText(s string, min, max int) bool {
	n := utf8.RuneCountInString(s)
	if n < min || n > max || strings.TrimSpace(s) == "" {
		return false
	}

	for _, r := range s {
		if r < 32 && r != '\n' && r != '\t' {
			return false
		}
	}

	return true
}

// ValidateRegistration проверяет формат данных и обязательное согласие до записи в БД.
func ValidateRegistration(in Registration) error {
	if !in.Consent {
		return ErrInvalid
	}

	return validateAccount(in)
}

// validateAccount проверяет общие поля аккаунта для сайта и административной CLI.
func validateAccount(in Registration) error {
	address, err := mail.ParseAddress(in.Email)
	if !loginPattern.MatchString(in.Login) ||
		!validText(in.Name, 2, 60) ||
		in.Name != strings.TrimSpace(in.Name) ||
		err != nil ||
		address.Address != in.Email ||
		len(in.Email) > 254 ||
		!strings.Contains(in.Email, ".") {
		return ErrInvalid
	}

	if utf8.RuneCountInString(in.Password) < 10 {
		return ErrPasswordShort
	}

	if len(in.Password) > 72 {
		return ErrPasswordLong
	}

	if strings.TrimSpace(in.Password) == "" {
		return ErrInvalid
	}

	return nil
}

// Register шифрует персональные поля и сохраняет пароль в виде солёного хеша bcrypt.
func (s *Service) Register(ctx context.Context, in Registration) (*model.Account, error) {
	if err := ValidateRegistration(in); err != nil {
		return nil, err
	}

	return s.register(ctx, in, false)
}

// RegisterAdmin создаёт администратора только через отдельное подключение CLI.
func (s *Service) RegisterAdmin(ctx context.Context, in Registration) (*model.Account, error) {
	if err := validateAccount(in); err != nil {
		return nil, err
	}

	return s.register(ctx, in, true)
}

// register шифрует и сохраняет данные, проверенные выбранным сценарием создания аккаунта.
func (s *Service) register(ctx context.Context, in Registration, admin bool) (*model.Account, error) {
	_, err := s.repo.ByLogin(ctx, s.Vault.Index(in.Login))
	if err == nil {
		return nil, ErrConflict
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	a := &model.Account{LoginIndex: s.Vault.Index(in.Login)}
	for value, target := range map[string]*string{
		in.Login: &a.LoginEnc,
		in.Email: &a.EmailEnc,
	} {
		*target, err = s.Vault.Encrypt(value)
		if err != nil {
			return nil, err
		}
	}

	a.NameEnc, err = s.Vault.Encrypt(in.Name)
	if err != nil {
		return nil, err
	}

	a.BioEnc, err = s.Vault.Encrypt("Пока ничего о себе не рассказал.")
	if err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), 12)
	if err != nil {
		return nil, err
	}

	a.PasswordHash = string(hash)
	if err = s.repo.CreateAccount(ctx, a, admin); err != nil {
		var pg *pq.Error
		if errors.As(err, &pg) && pg.Code == "23505" {
			return nil, ErrConflict
		}

		return nil, err
	}

	return s.repo.Account(ctx, a.ID)
}

// Login возвращает одинаковую ошибку для неизвестного логина, пароля и блокировки.
func (s *Service) Login(ctx context.Context, login, password string) (*model.Account, error) {
	if !loginPattern.MatchString(login) || len(password) < 10 || len(password) > 72 {
		return nil, ErrCredentials
	}

	a, err := s.repo.ByLogin(ctx, s.Vault.Index(login))
	if errors.Is(err, sql.ErrNoRows) {
		_ = bcrypt.CompareHashAndPassword(s.dummyHash, []byte(password))
		return nil, ErrCredentials
	}

	if err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(password)) != nil || a.Banned {
		return nil, ErrCredentials
	}

	return a, nil
}

// Public раскрывает только публичные поля; почта доступна владельцу аккаунта.
func (s *Service) Public(a *model.Account, private bool) (*model.User, error) {
	u := &model.User{
		ID:        a.ID,
		Role:      a.Role,
		Banned:    a.Banned,
		CreatedAt: a.CreatedAt,
	}
	var err error
	for encrypted, target := range map[string]*string{
		a.LoginEnc: &u.Login,
		a.NameEnc:  &u.Name,
		a.BioEnc:   &u.Bio,
	} {
		*target, err = s.Vault.Decrypt(encrypted)
		if err != nil {
			return nil, err
		}
	}

	if private {
		u.Email, err = s.Vault.Decrypt(a.EmailEnc)
	}

	return u, err
}

// User загружает публичный профиль по идентификатору.
func (s *Service) User(ctx context.Context, id int64) (*model.User, error) {
	a, err := s.repo.Account(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.Public(a, false)
}

// UpdateProfile проверяет и шифрует редактируемые поля профиля.
func (s *Service) UpdateProfile(ctx context.Context, id int64, name, bio string) error {
	bio = normalizeText(bio)
	if !validText(name, 2, 60) || !validText(bio, 1, 500) {
		return ErrInvalid
	}

	n, err := s.Vault.Encrypt(strings.TrimSpace(name))
	if err != nil {
		return err
	}

	b, err := s.Vault.Encrypt(strings.TrimSpace(bio))
	if err != nil {
		return err
	}

	return s.repo.Profile(ctx, id, n, b)
}

// Posts добавляет публичные данные авторов к странице публикаций.
func (s *Service) Posts(ctx context.Context, author int64, offset uint64) ([]model.Post, error) {
	out, err := s.repo.Posts(ctx, author, offset)
	if err != nil {
		return nil, err
	}

	for i := range out {
		out[i].Author, err = s.User(ctx, out[i].AuthorID)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

// Post загружает публикацию и её автора.
func (s *Service) Post(ctx context.Context, id int64) (*model.Post, error) {
	p, err := s.repo.Post(ctx, id)
	if err != nil {
		return nil, err
	}

	p.Author, err = s.User(ctx, p.AuthorID)
	return p, err
}

// CreatePost проверяет ограничения заголовка и текста публикации.
func (s *Service) CreatePost(ctx context.Context, author int64, title, body string) (int64, error) {
	title = strings.TrimSpace(title)
	body = strings.TrimSpace(normalizeText(body))
	if !validText(title, 3, 120) || !validText(body, 1, 10000) {
		return 0, ErrInvalid
	}

	return s.repo.CreatePost(ctx, author, title, body)
}

// DeletePost разрешает удаление только автору или администратору.
func (s *Service) DeletePost(ctx context.Context, a *model.Account, id int64) error {
	p, err := s.repo.Post(ctx, id)
	if err != nil {
		return err
	}

	if a.Role != "admin" && p.AuthorID != a.ID {
		return ErrForbidden
	}

	return s.repo.DeletePost(ctx, id)
}

// Comments раскрывает публичных авторов комментариев.
func (s *Service) Comments(ctx context.Context, id int64) ([]model.Comment, error) {
	if _, err := s.repo.Post(ctx, id); err != nil {
		return nil, err
	}

	out, err := s.repo.Comments(ctx, id)
	if err != nil {
		return nil, err
	}

	for i := range out {
		out[i].Author, err = s.User(ctx, out[i].AuthorID)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

// Comment проверяет текст и существование публикации до добавления комментария.
func (s *Service) Comment(ctx context.Context, post, author int64, body string) error {
	body = strings.TrimSpace(normalizeText(body))
	if !validText(body, 1, 2000) {
		return ErrInvalid
	}

	if _, err := s.repo.Post(ctx, post); err != nil {
		return err
	}

	err := s.repo.CreateComment(ctx, post, author, body)
	var pg *pq.Error
	if errors.As(err, &pg) && pg.Code == "23503" {
		return ErrNotFound
	}

	return err
}

// Users возвращает список публичных профилей для модерации.
func (s *Service) Users(ctx context.Context) ([]*model.User, error) {
	accounts, err := s.repo.Users(ctx)
	if err != nil {
		return nil, err
	}

	out := []*model.User{}
	for i := range accounts {
		u, err := s.Public(&accounts[i], false)
		if err != nil {
			return nil, err
		}

		out = append(out, u)
	}

	return out, nil
}
