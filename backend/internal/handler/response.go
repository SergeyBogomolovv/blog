// Строгий JSON и публичные сообщения об ошибках.
package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"mime"
	"net/http"

	"cleanhome/internal/service"
)

// errorResponse безопасное сообщение и ошибки отдельных полей.
type errorResponse struct {
	Error  string         `json:"error"`
	Fields service.Fields `json:"fields"`
}

// respond записывает JSON-ответ с заданным HTTP-статусом.
func respond(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// fail преобразует ошибки сервиса в безопасные сообщения клиенту.
func (h *Handler) fail(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	message := "Не удалось выполнить запрос. Попробуйте позже."
	fields := service.Fields{}

	switch {
	case errors.As(err, &fields):
		status = http.StatusBadRequest
		message = fields.Error()
	case errors.Is(err, service.ErrInvalid):
		status = http.StatusBadRequest
		message = err.Error()
	case errors.Is(err, service.ErrCredentials), errors.Is(err, service.ErrUnauthorized):
		status = http.StatusUnauthorized
		message = err.Error()
	case errors.Is(err, service.ErrForbidden):
		status = http.StatusForbidden
		message = err.Error()
	case errors.Is(err, service.ErrConflict):
		status = http.StatusConflict
		message = err.Error()
		fields["login"] = message
	case errors.Is(err, service.ErrNotFound):
		status = http.StatusNotFound
		message = err.Error()
	}

	if status == http.StatusInternalServerError {
		slog.Error("Ошибка обработки запроса", "error", err)
	}

	respond(w, status, errorResponse{
		Error:  message,
		Fields: fields,
	})
}

// decode принимает один JSON-объект и отклоняет неизвестные поля.
func decode(w http.ResponseWriter, r *http.Request, target any) error {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		return service.ErrInvalid
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBody)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(target); err != nil {
		return service.ErrInvalid
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return service.ErrInvalid
	}

	return nil
}
