// Защита запросов и проверка сохраняемых в PostgreSQL сессий.
package handler

import (
	"bytes"
	"context"
	"crypto/hmac"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"

	"cleanhome/internal/model"
	"cleanhome/internal/service"
	"cleanhome/pkg/security"
)

const apiContentSecurityPolicy = "default-src 'self'; script-src 'self'; " +
	"object-src 'none'; base-uri 'self'; form-action 'self'; frame-ancestors 'none'"

// security добавляет защитные заголовки и проверяет изменяющие запросы.
func (h *Handler) security(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		w.Header().Set("Content-Security-Policy", apiContentSecurityPolicy)

		if !isReadOnlyMethod(r.Method) {
			if !h.limitBody(w, r) {
				return
			}
			if !h.validCSRFRequest(r) {
				h.fail(w, service.ErrForbidden)
				return
			}
		}

		defer func() {
			if recover() != nil {
				slog.Error("Непредвиденная ошибка обработчика")
				respond(w, http.StatusInternalServerError, map[string]string{
					"error": "Внутренняя ошибка сервера",
				})
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// isReadOnlyMethod определяет методы, не изменяющие данные.
func isReadOnlyMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
}

// limitBody ограничивает размер тела запроса до обработки JSON.
func (h *Handler) limitBody(w http.ResponseWriter, r *http.Request) bool {
	// Ограничение размера действует и для выхода, и для подтверждения cookie.
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBody)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			respond(w, http.StatusRequestEntityTooLarge, map[string]string{
				"error": "Запрос превышает 16 КБ.",
			})
		} else {
			h.fail(w, service.ErrInvalid)
		}

		return false
	}

	r.Body = io.NopCloser(bytes.NewReader(body))
	return true
}

// validCSRFRequest проверяет Origin, служебный заголовок и подпись CSRF.
func (h *Handler) validCSRFRequest(r *http.Request) bool {
	if r.Header.Get("Origin") != h.config.Origin {
		return false
	}
	if r.Header.Get("X-Requested-With") != "cleaning" {
		return false
	}

	cookie, err := r.Cookie(csrfCookieName)
	if err != nil {
		return false
	}

	token := r.Header.Get("X-CSRF-Token")
	if !hmac.Equal([]byte(cookie.Value), []byte(token)) {
		return false
	}

	return security.ValidCSRF(h.config.SessionKey, token, sessionValue(r))
}

// setCookie устанавливает срок действия и защитные параметры cookie.
func setCookie(w http.ResponseWriter, name, value string, maxAge int) {
	expires := time.Now().Add(time.Duration(maxAge) * time.Second)
	if maxAge < 0 {
		expires = time.Unix(1, 0)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
		Expires:  expires,
	})
}

// sessionValue читает токен сессии из cookie запроса.
func sessionValue(r *http.Request) string {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return ""
	}

	return cookie.Value
}

// session проверяет подпись, срок и наличие сессии для текущего запроса.
func (h *Handler) session(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := sessionValue(r)
		isAuthRequest := r.URL.Path == "/api/auth/login" || r.URL.Path == "/api/auth/register"
		if token == "" || isAuthRequest {
			next.ServeHTTP(w, r)
			return
		}
		if !security.ValidSession(h.config.SessionKey, token) {
			setCookie(w, sessionCookieName, "", -1)
			next.ServeHTTP(w, r)
			return
		}

		user, err := h.service.Authorize(r.Context(), security.SessionHash(token))
		if errors.Is(err, service.ErrUnauthorized) {
			setCookie(w, sessionCookieName, "", -1)
			next.ServeHTTP(w, r)
			return
		}
		if err != nil {
			h.fail(w, err)
			return
		}

		setCookie(w, sessionCookieName, token, sessionMaxAge)
		ctx := context.WithValue(r.Context(), contextKey{}, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// require запрещает защищённую операцию без авторизованного пользователя.
func (h *Handler) require(w http.ResponseWriter, r *http.Request) *model.User {
	user, _ := r.Context().Value(contextKey{}).(*model.User)
	if user == nil {
		h.fail(w, service.ErrUnauthorized)
	}

	return user
}

// Next.js и прямые вызовы API разделяют лимит; X-Forwarded-For не используется.
func (h *Handler) limited(w http.ResponseWriter, r *http.Request) bool {
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	now := time.Now()

	h.mu.Lock()
	defer h.mu.Unlock()

	for address, previous := range h.attempts {
		if !now.Before(previous.until) {
			delete(h.attempts, address)
		}
	}

	current := h.attempts[host]
	if current.until.IsZero() {
		current.until = now.Add(authAttemptWindow)
	}
	current.count++
	h.attempts[host] = current

	if current.count <= maxAuthAttempts {
		return false
	}

	retryAfter := int(time.Until(current.until).Seconds()) + 1
	w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
	respond(w, http.StatusTooManyRequests, map[string]string{
		"error": "Слишком много попыток. Повторите через 15 минут.",
	})

	return true
}
