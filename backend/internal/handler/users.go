// Обработчики профилей и административной модерации.
package handler

import (
	"net/http"

	"studyblog/internal/service"
)

// user возвращает публичную страницу пользователя без его почты.
func (h *Handler) user(w http.ResponseWriter, r *http.Request) {
	n, err := id(r)
	if err != nil {
		h.fail(w, err)
		return
	}

	u, err := h.service.User(r.Context(), n)
	if err != nil {
		h.fail(w, err)
		return
	}

	respond(w, 200, u)
}

// profile сохраняет только профиль владельца сессии.
func (h *Handler) profile(w http.ResponseWriter, r *http.Request) {
	a := h.require(w, r, false)
	if a == nil {
		return
	}

	var in struct {
		Name string `json:"name"`
		Bio  string `json:"bio"`
	}

	err := decode(w, r, &in)
	if err == nil {
		err = h.service.UpdateProfile(r.Context(), a.ID, in.Name, in.Bio)
	}

	if err != nil {
		h.fail(w, err)
		return
	}

	respond(w, 200, map[string]bool{"ok": true})
}

// users возвращает список пользователей только администратору.
func (h *Handler) users(w http.ResponseWriter, r *http.Request) {
	if h.require(w, r, true) == nil {
		return
	}

	out, err := h.service.Users(r.Context())
	if err != nil {
		h.fail(w, err)
		return
	}

	respond(w, 200, out)
}

// ban блокирует или разблокирует обычный аккаунт.
func (h *Handler) ban(w http.ResponseWriter, r *http.Request) {
	if h.require(w, r, true) == nil {
		return
	}

	n, err := id(r)
	if err != nil {
		h.fail(w, err)
		return
	}

	var in struct {
		Banned *bool `json:"banned"`
	}

	err = decode(w, r, &in)
	if err == nil && in.Banned == nil {
		err = service.ErrInvalid
	}

	if err == nil {
		err = h.service.Ban(r.Context(), current(r), n, *in.Banned)
	}

	if err != nil {
		h.fail(w, err)
		return
	}

	respond(w, 200, map[string]bool{"ok": true})
}
