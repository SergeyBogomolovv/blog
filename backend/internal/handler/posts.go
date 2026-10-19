// Обработчики публикаций и комментариев.
package handler

import (
	"net/http"
	"strconv"

	"studyblog/internal/service"
)

// posts возвращает страницу ленты с необязательным фильтром по автору.
func (h *Handler) posts(w http.ResponseWriter, r *http.Request) {
	var author int64
	var page uint64
	var err error
	if v := r.URL.Query().Get("author"); v != "" {
		author, err = strconv.ParseInt(v, 10, 64)
		if err != nil || author < 1 {
			h.fail(w, service.ErrInvalid)
			return
		}
	}

	if v := r.URL.Query().Get("page"); v != "" {
		page, err = strconv.ParseUint(v, 10, 32)
		if err != nil || page > 10000 {
			h.fail(w, service.ErrInvalid)
			return
		}
	}

	out, err := h.service.Posts(r.Context(), author, page*20)
	if err != nil {
		h.fail(w, err)
		return
	}

	respond(w, 200, out)
}

// post возвращает отдельную публикацию.
func (h *Handler) post(w http.ResponseWriter, r *http.Request) {
	n, err := id(r)
	if err != nil {
		h.fail(w, err)
		return
	}

	p, err := h.service.Post(r.Context(), n)
	if err != nil {
		h.fail(w, err)
		return
	}

	respond(w, 200, p)
}

// createPost публикует запись от имени текущего пользователя.
func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	a := h.require(w, r, false)
	if a == nil {
		return
	}

	var in struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	if err := decode(w, r, &in); err != nil {
		h.fail(w, err)
		return
	}

	n, err := h.service.CreatePost(r.Context(), a.ID, in.Title, in.Body)
	if err != nil {
		h.fail(w, err)
		return
	}

	respond(w, 201, map[string]int64{"id": n})
}

// deletePost удаляет запись после проверки прав сервисом.
func (h *Handler) deletePost(w http.ResponseWriter, r *http.Request) {
	a := h.require(w, r, false)
	if a == nil {
		return
	}

	n, err := id(r)
	if err == nil {
		err = h.service.DeletePost(r.Context(), a, n)
	}

	if err != nil {
		h.fail(w, err)
		return
	}

	respond(w, 200, map[string]bool{"ok": true})
}

// comments возвращает комментарии к существующей публикации.
func (h *Handler) comments(w http.ResponseWriter, r *http.Request) {
	n, err := id(r)
	if err != nil {
		h.fail(w, err)
		return
	}

	out, err := h.service.Comments(r.Context(), n)
	if err != nil {
		h.fail(w, err)
		return
	}

	respond(w, 200, out)
}

// comment добавляет комментарий авторизованного пользователя.
func (h *Handler) comment(w http.ResponseWriter, r *http.Request) {
	a := h.require(w, r, false)
	if a == nil {
		return
	}

	n, err := id(r)
	if err != nil {
		h.fail(w, err)
		return
	}

	var in struct {
		Body string `json:"body"`
	}

	if err = decode(w, r, &in); err == nil {
		err = h.service.Comment(r.Context(), n, a.ID, in.Body)
	}

	if err != nil {
		h.fail(w, err)
		return
	}

	respond(w, 201, map[string]bool{"ok": true})
}
