// Сериализация JSON, строгий ввод и единые ответы об ошибках.
package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"studyblog/internal/service"
)

// respond сериализует JSON без вывода внутренних данных ошибок.
func respond(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

// fail сопоставляет доменные ошибки с кодами HTTP.
func (h *Handler) fail(w http.ResponseWriter, err error) {
	status, message := 500, "Не удалось выполнить запрос. Попробуйте позже."
	switch {
	case errors.Is(err, service.ErrInvalid),
		errors.Is(err, service.ErrPasswordShort),
		errors.Is(err, service.ErrPasswordLong):
		status, message = 400, err.Error()
	case errors.Is(err, service.ErrCredentials), errors.Is(err, service.ErrUnauthorized):
		status, message = 401, err.Error()
	case errors.Is(err, service.ErrForbidden):
		status, message = 403, err.Error()
	case errors.Is(err, service.ErrConflict):
		status, message = 409, err.Error()
	case errors.Is(err, sql.ErrNoRows), errors.Is(err, service.ErrNotFound):
		status, message = 404, service.ErrNotFound.Error()
	}

	if status == 500 {
		slog.Error("Ошибка обработки запроса", "error", err)
	}

	respond(w, status, map[string]string{"error": message})
}

// decode ограничивает тело запроса и запрещает неизвестные поля и лишний JSON.
func decode(w http.ResponseWriter, r *http.Request, target any) error {
	kind, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || kind != "application/json" {
		return service.ErrInvalid
	}

	r.Body = http.MaxBytesReader(w, r.Body, 65536)
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if d.Decode(target) != nil {
		return service.ErrInvalid
	}

	if d.Decode(&struct{}{}) != io.EOF {
		return service.ErrInvalid
	}

	return nil
}

// id извлекает положительный идентификатор из маршрута.
func id(r *http.Request) (int64, error) {
	n, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || n < 1 {
		return 0, service.ErrInvalid
	}

	return n, nil
}
