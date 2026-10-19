// Обработчики регистрации, входа, выхода и состояния сервера.
package handler

import (
	"net/http"

	"studyblog/internal/model"
	"studyblog/internal/service"
	"studyblog/pkg/security"
)

// health проверяет доступность базы данных.
func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Health(r.Context()); err != nil {
		h.fail(w, err)
		return
	}

	respond(w, 200, map[string]string{"status": "ok"})
}

// me возвращает текущего пользователя или null для гостя.
func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	a := current(r)
	if a == nil {
		respond(w, 200, map[string]any{"user": nil})
		return
	}

	u, err := h.service.Public(a, true)
	if err != nil {
		h.fail(w, err)
		return
	}

	respond(w, 200, map[string]any{"user": u})
}

// authenticate выдаёт access-cookie и публичный результат входа.
func (h *Handler) authenticate(w http.ResponseWriter, a *model.Account, status int) {
	token, err := security.Issue(h.config.JWTKey, a.ID, a.TokenVersion)
	if err != nil {
		h.fail(w, err)
		return
	}

	u, err := h.service.Public(a, true)
	if err != nil {
		h.fail(w, err)
		return
	}

	h.cookie(w, token, false)
	respond(w, status, map[string]any{"user": u})
}

// register принимает данные только с явно установленным согласием.
func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	if h.limited(w, r) {
		return
	}

	var in service.Registration
	if err := decode(w, r, &in); err != nil {
		h.fail(w, err)
		return
	}

	a, err := h.service.Register(r.Context(), in)
	if err != nil {
		h.fail(w, err)
		return
	}

	h.authenticate(w, a, 201)
}

// login проверяет логин и пароль без раскрытия причины отказа.
func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	if h.limited(w, r) {
		return
	}

	var in struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := decode(w, r, &in); err != nil {
		h.fail(w, err)
		return
	}

	a, err := h.service.Login(r.Context(), in.Login, in.Password)
	if err != nil {
		h.fail(w, err)
		return
	}

	h.authenticate(w, a, 200)
}

// logout отзывает все токены пользователя и удаляет cookie.
func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	if a := current(r); a != nil {
		if err := h.service.Logout(r.Context(), a.ID); err != nil {
			h.fail(w, err)
			return
		}
	}

	h.cookie(w, "", true)
	respond(w, 200, map[string]bool{"ok": true})
}
