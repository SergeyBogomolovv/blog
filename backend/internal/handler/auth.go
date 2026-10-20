package handler

import (
	"net/http"
	"strings"

	"cleanhome/internal/service"
	"cleanhome/pkg/security"
)

type loginInput struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	if h.limited(w, r) {
		return
	}

	var input service.Registration
	if err := decode(w, r, &input); err != nil {
		h.fail(w, err)
		return
	}
	if err := h.service.Register(r.Context(), input); err != nil {
		h.fail(w, err)
		return
	}

	respond(w, http.StatusCreated, map[string]bool{"ok": true})
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	if h.limited(w, r) {
		return
	}

	var input loginInput
	if err := decode(w, r, &input); err != nil {
		h.fail(w, err)
		return
	}

	user, err := h.service.Login(r.Context(), strings.TrimSpace(input.Login), input.Password)
	if err != nil {
		h.fail(w, err)
		return
	}

	token, err := security.SessionID(h.config.SessionKey)
	if err != nil {
		h.fail(w, err)
		return
	}
	if err := h.revokeCurrentSession(r); err != nil {
		h.fail(w, err)
		return
	}

	sessionID := security.SessionHash(token)
	if err := h.service.CreateSession(r.Context(), sessionID, user.ID); err != nil {
		h.fail(w, err)
		return
	}

	setCookie(w, sessionCookieName, token, sessionMaxAge)
	setCookie(w, csrfCookieName, "", -1)
	respond(w, http.StatusOK, map[string]any{"user": user})
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	user := h.require(w, r)
	if user == nil {
		return
	}

	respond(w, http.StatusOK, map[string]any{"user": user})
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	// Подтверждение выхода отправляется после удаления сессии из базы.
	if err := h.revokeCurrentSession(r); err != nil {
		h.fail(w, err)
		return
	}

	setCookie(w, sessionCookieName, "", -1)
	setCookie(w, csrfCookieName, "", -1)
	respond(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *Handler) revokeCurrentSession(r *http.Request) error {
	token := sessionValue(r)
	if !security.ValidSession(h.config.SessionKey, token) {
		return nil
	}

	return h.service.Logout(r.Context(), security.SessionHash(token))
}
