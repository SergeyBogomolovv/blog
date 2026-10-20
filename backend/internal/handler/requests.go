package handler

import (
	"net/http"

	"cleanhome/internal/service"
)

func (h *Handler) requests(w http.ResponseWriter, r *http.Request) {
	user := h.require(w, r)
	if user == nil {
		return
	}

	requests, err := h.service.Requests(r.Context(), user.ID)
	if err != nil {
		h.fail(w, err)
		return
	}

	respond(w, http.StatusOK, map[string]any{"requests": requests})
}

func (h *Handler) createRequest(w http.ResponseWriter, r *http.Request) {
	user := h.require(w, r)
	if user == nil {
		return
	}

	var input service.RequestInput
	if err := decode(w, r, &input); err != nil {
		h.fail(w, err)
		return
	}

	request, err := h.service.CreateRequest(r.Context(), user.ID, input)
	if err != nil {
		h.fail(w, err)
		return
	}

	respond(w, http.StatusCreated, map[string]any{"request": request})
}
