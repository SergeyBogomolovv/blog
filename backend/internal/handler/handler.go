// Package handler предоставляет JSON API клинингового сервиса.
package handler

import (
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	"cleanhome/internal/config"
	"cleanhome/internal/service"
	"cleanhome/pkg/security"
)

const (
	sessionCookieName  = "cleaning.sid"
	csrfCookieName     = "csrf_token"
	sessionMaxAge      = 30 * 60
	cookieNoticeMaxAge = 180 * 24 * 60 * 60
	maxRequestBody     = 16 * 1024
	maxAuthAttempts    = 30
	authAttemptWindow  = 15 * time.Minute
)

// contextKey закрытый ключ пользователя в контексте HTTP-запроса.
type contextKey struct{}

// attempt счётчик и срок окна ограничения попыток.
type attempt struct {
	count int
	until time.Time
}

// Handler обработчики API и состояние ограничения попыток входа.
type Handler struct {
	service  *service.Service
	config   config.Config
	mu       sync.Mutex
	attempts map[string]attempt
}

// New регистрирует маршруты API и общие проверки безопасности.
func New(svc *service.Service, cfg config.Config) http.Handler {
	h := &Handler{
		service:  svc,
		config:   cfg,
		attempts: make(map[string]attempt),
	}

	router := chi.NewRouter()
	router.Use(h.security, h.session)

	router.Get("/api/health", h.health)
	router.Get("/api/services", h.services)
	router.Get("/api/csrf", h.csrf)
	router.Post("/api/cookie-notice", h.acceptCookieNotice)

	router.Post("/api/auth/register", h.register)
	router.Post("/api/auth/login", h.login)
	router.Post("/api/auth/logout", h.logout)
	router.Get("/api/auth/me", h.me)

	router.Get("/api/requests", h.requests)
	router.Post("/api/requests", h.createRequest)
	router.NotFound(h.notFound)

	return router
}

// health проверяет доступность сервиса и базы данных.
func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Health(r.Context()); err != nil {
		h.fail(w, err)
		return
	}

	respond(w, http.StatusOK, map[string]bool{"ok": true})
}

// services возвращает каталог услуг с учётом поискового запроса.
func (h *Handler) services(w http.ResponseWriter, r *http.Request) {
	services, err := service.Search(r.URL.Query().Get("q"))
	if err != nil {
		h.fail(w, err)
		return
	}

	respond(w, http.StatusOK, map[string]any{"services": services})
}

// csrf выдаёт подписанный токен для отправки формы.
func (h *Handler) csrf(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(csrfCookieName)
	if err == nil && security.ValidCSRF(h.config.SessionKey, cookie.Value, sessionValue(r)) {
		// Параллельные формы могут продолжать пользоваться тем же токеном.
		setCookie(w, csrfCookieName, cookie.Value, sessionMaxAge)
		respond(w, http.StatusOK, map[string]string{"token": cookie.Value})
		return
	}

	token, err := security.NewCSRF(h.config.SessionKey, sessionValue(r))
	if err != nil {
		h.fail(w, err)
		return
	}

	setCookie(w, csrfCookieName, token, sessionMaxAge)
	respond(w, http.StatusOK, map[string]string{"token": token})
}

// acceptCookieNotice запоминает подтверждение уведомления на 180 дней.
func (h *Handler) acceptCookieNotice(w http.ResponseWriter, _ *http.Request) {
	setCookie(w, "cookie_notice", "accepted", cookieNoticeMaxAge)
	respond(w, http.StatusOK, map[string]bool{"ok": true})
}

// notFound возвращает JSON-ответ для отсутствующего маршрута.
func (h *Handler) notFound(w http.ResponseWriter, _ *http.Request) {
	h.fail(w, service.ErrNotFound)
}
