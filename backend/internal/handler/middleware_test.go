// Проверки защитных заголовков, CSRF и ограничений запросов.
package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cleanhome/internal/config"
	"cleanhome/pkg/security"
)

// TestSecurity проверяет CSRF-защиту и обязательные заголовки ответов.
func TestSecurity(t *testing.T) {
	key := bytes.Repeat([]byte{1}, 32)
	h := &Handler{
		config: config.Config{
			Origin:     "http://127.0.0.1:3107",
			SessionKey: key,
		},
	}
	token, err := security.NewCSRF(key, "")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		origin string
		cookie string
		header string
		status int
	}{
		{
			name:   "valid",
			origin: h.config.Origin,
			cookie: token,
			header: token,
			status: http.StatusOK,
		},
		{
			name:   "no csrf",
			origin: h.config.Origin,
			status: http.StatusForbidden,
		},
		{
			name:   "forged csrf",
			origin: h.config.Origin,
			cookie: "fake",
			header: "fake",
			status: http.StatusForbidden,
		},
		{
			name:   "different token",
			origin: h.config.Origin,
			cookie: token,
			header: token + "x",
			status: http.StatusForbidden,
		},
		{
			name:   "foreign origin",
			origin: "http://evil.test",
			cookie: token,
			header: token,
			status: http.StatusForbidden,
		},
	}
	headerNames := []string{
		"Content-Security-Policy",
		"X-Frame-Options",
		"X-Content-Type-Options",
		"Permissions-Policy",
		"Referrer-Policy",
		"Cache-Control",
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/requests", strings.NewReader(`{}`))
			request.Header.Set("Origin", tc.origin)
			request.Header.Set("X-Requested-With", "cleaning")
			request.Header.Set("X-CSRF-Token", tc.header)
			request.AddCookie(&http.Cookie{Name: csrfCookieName, Value: tc.cookie})

			next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			response := httptest.NewRecorder()
			h.security(next).ServeHTTP(response, request)

			if response.Code != tc.status {
				t.Fatalf("status = %d, want %d", response.Code, tc.status)
			}
			for _, name := range headerNames {
				if response.Header().Get(name) == "" {
					t.Errorf("missing response header: %s", name)
				}
			}
		})
	}
}

// TestJSONLimit проверяет ограничения размера и структуры JSON.
func TestJSONLimit(t *testing.T) {
	bodies := []string{
		`{"unknown":1}`,
		`{} {}`,
		`{"date":"` + strings.Repeat("a", 16*1024) + `"}`,
	}
	for _, body := range bodies {
		request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		var input struct {
			Date string `json:"date"`
		}
		if err := decode(response, request, &input); err == nil {
			t.Fatal("invalid body accepted")
		}
	}
}

// TestRateLimit проверяет блокировку после превышения числа попыток.
func TestRateLimit(t *testing.T) {
	h := &Handler{attempts: make(map[string]attempt)}
	request := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)

	for number := 1; number <= 31; number++ {
		response := httptest.NewRecorder()
		limited := h.limited(response, request)
		wantLimited := number > 30

		if limited != wantLimited {
			t.Fatalf("request %d: limited = %v, want %v", number, limited, wantLimited)
		}
		if limited && response.Code != http.StatusTooManyRequests {
			t.Fatalf("status = %d, want 429", response.Code)
		}
	}
}

// TestAllMutationBodiesLimited проверяет лимит тела для всех изменяющих маршрутов.
func TestAllMutationBodiesLimited(t *testing.T) {
	h := &Handler{}
	body := strings.NewReader(strings.Repeat("x", 16*1024+1))
	request := httptest.NewRequest(http.MethodPost, "/api/cookie-notice", body)
	response := httptest.NewRecorder()
	next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		t.Fatal("oversized request reached handler")
	})
	h.security(next).ServeHTTP(response, request)

	if response.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", response.Code)
	}
}
