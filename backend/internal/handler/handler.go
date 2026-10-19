// Маршруты HTTP API и зависимости обработчиков.
package handler

import (
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	"studyblog/internal/config"
	"studyblog/internal/service"
)

type contextKey struct{}
type attempt struct {
	count int
	until time.Time
}

type Handler struct {
	service  *service.Service
	config   config.Config
	mu       sync.Mutex
	attempts map[string]attempt
}

// New собирает маршруты, общую проверку сессии и защиту запросов изменения данных.
func New(s *service.Service, c config.Config) http.Handler {
	h := &Handler{service: s, config: c, attempts: make(map[string]attempt)}
	r := chi.NewRouter()
	r.Use(h.security, h.session)

	r.Get("/api/health", h.health)

	r.Get("/api/auth/me", h.me)
	r.Post("/api/auth/register", h.register)
	r.Post("/api/auth/login", h.login)
	r.Post("/api/auth/logout", h.logout)

	r.Get("/api/posts", h.posts)
	r.Get("/api/posts/{id}", h.post)
	r.Get("/api/posts/{id}/comments", h.comments)
	r.Post("/api/posts", h.createPost)
	r.Delete("/api/posts/{id}", h.deletePost)
	r.Post("/api/posts/{id}/comments", h.comment)

	r.Get("/api/users/{id}", h.user)
	r.Patch("/api/profile", h.profile)

	r.Get("/api/admin/users", h.users)
	r.Patch("/api/admin/users/{id}/ban", h.ban)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		h.fail(w, service.ErrNotFound)
	})
	return r
}
