// Проверка сессий, защита от CSRF и ограничение частоты входа.
package handler

import (
	"context"
	"crypto/hmac"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"

	"studyblog/internal/model"
	"studyblog/internal/service"
	"studyblog/pkg/security"
)

// security проверяет Origin, специальный заголовок и подписанный CSRF-токен.
func (h *Handler) security(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		if r.Method != "GET" && r.Method != "HEAD" && r.Method != "OPTIONS" {
			if r.Header.Get("Origin") != h.config.Origin || r.Header.Get("X-Requested-With") != "studyblog" {
				h.fail(w, service.ErrForbidden)
				return
			}
			csrf, err := r.Cookie("csrf_token")
			session := ""
			if cookie, err := r.Cookie("access_token"); err == nil {
				session = cookie.Value
			}
			provided := r.Header.Get("X-CSRF-Token")
			if err != nil || !hmac.Equal([]byte(csrf.Value), []byte(provided)) ||
				!security.ValidCSRF(h.config.JWTKey, provided, session) {
				h.fail(w, service.ErrForbidden)
				return
			}
		}

		defer func() {
			if recover() != nil {
				slog.Error("Непредвиденная ошибка обработчика")
				respond(w, 500, map[string]string{"error": "Внутренняя ошибка сервера"})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// cookie устанавливает HttpOnly-cookie с единственным access-токеном.
func (h *Handler) cookie(w http.ResponseWriter, value string, remove bool) {
	age := 7200
	expires := time.Now().Add(2 * time.Hour)
	if remove {
		age = -1
		expires = time.Unix(1, 0)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.config.Secure,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   age,
		Expires:  expires,
	})
}

// session на каждом запросе проверяет подпись, срок, блокировку и версию сессии.
func (h *Handler) session(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Вход и регистрация доступны с истёкшей cookie; CSRF уже проверен выше.
		if r.Method == http.MethodPost && (r.URL.Path == "/api/auth/login" || r.URL.Path == "/api/auth/register") {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("access_token")
		if errors.Is(err, http.ErrNoCookie) {
			next.ServeHTTP(w, r)
			return
		}

		if err != nil {
			h.fail(w, service.ErrUnauthorized)
			return
		}

		c, err := security.Parse(h.config.JWTKey, cookie.Value)
		if err != nil {
			h.cookie(w, "", true)
			h.fail(w, service.ErrUnauthorized)
			return
		}

		id, err := strconv.ParseInt(c.Subject, 10, 64)
		if err != nil {
			h.cookie(w, "", true)
			h.fail(w, service.ErrUnauthorized)
			return
		}

		a, err := h.service.Authorize(r.Context(), id, c.Version)
		if err != nil {
			if errors.Is(err, service.ErrUnauthorized) {
				h.cookie(w, "", true)
			}

			h.fail(w, err)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKey{}, a)))
	})
}

// current получает проверенную middleware учётную запись.
func current(r *http.Request) *model.Account {
	a, _ := r.Context().Value(contextKey{}).(*model.Account)
	return a
}

// require запрещает действия гостям и проверяет административную роль.
func (h *Handler) require(w http.ResponseWriter, r *http.Request, admin bool) *model.Account {
	a := current(r)
	if a == nil {
		h.fail(w, service.ErrUnauthorized)
		return nil
	}

	if admin && a.Role != "admin" {
		h.fail(w, service.ErrForbidden)
		return nil
	}

	return a
}

// limited ограничивает попытки входа и регистрации в течение минуты.
func (h *Handler) limited(w http.ResponseWriter, r *http.Request) bool {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	h.mu.Lock()
	defer h.mu.Unlock()
	now := time.Now()
	for key, a := range h.attempts {
		if now.After(a.until) {
			delete(h.attempts, key)
		}
	}

	a := h.attempts[host]
	if a.until.IsZero() {
		a.until = now.Add(time.Minute)
	}

	a.count++
	h.attempts[host] = a
	if a.count > 15 {
		w.Header().Set("Retry-After", "60")
		respond(w, 429, map[string]string{"error": "Слишком много попыток. Подождите минуту."})
		return true
	}

	return false
}
