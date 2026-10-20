package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"studyblog/internal/config"
)

const guestCSRF = "0000000000000000000000000000000000000000000000000000000000000000.60ca8f329f50bc54780ee5b29dc67301fc44cb31ff0eab9fee428bf3a178c548"
const sessionCSRF = "0000000000000000000000000000000000000000000000000000000000000000.2f897b9659eb48931898af16a1c0abbc3e58d2f5710d194a86463ee26e7f463c"

func TestSecurityCSRF(t *testing.T) {
	t.Parallel()
	h := &Handler{config: config.Config{Origin: "http://localhost:3000", JWTKey: bytes.Repeat([]byte{1}, 32)}}
	for _, tc := range []struct {
		name, method, origin, requestedWith, cookie, header, session string
		want                                                         int
	}{
		{"guest", "POST", h.config.Origin, "studyblog", guestCSRF, guestCSRF, "", 204},
		{"session", "POST", h.config.Origin, "studyblog", sessionCSRF, sessionCSRF, "session-a", 204},
		{"patch", "PATCH", h.config.Origin, "studyblog", guestCSRF, guestCSRF, "", 204},
		{"delete", "DELETE", h.config.Origin, "studyblog", guestCSRF, guestCSRF, "", 204},
		{"put", "PUT", h.config.Origin, "studyblog", guestCSRF, guestCSRF, "", 204},
		{"missing_both", "POST", h.config.Origin, "studyblog", "", "", "", 403},
		{"missing_cookie", "POST", h.config.Origin, "studyblog", "", guestCSRF, "", 403},
		{"missing_header", "POST", h.config.Origin, "studyblog", guestCSRF, "", "", 403},
		{"mismatched", "POST", h.config.Origin, "studyblog", guestCSRF, sessionCSRF, "", 403},
		{"forged_both", "POST", h.config.Origin, "studyblog", strings.Repeat("a", 129), strings.Repeat("a", 129), "", 403},
		{"another_session", "POST", h.config.Origin, "studyblog", sessionCSRF, sessionCSRF, "session-b", 403},
		{"guest_in_session", "POST", h.config.Origin, "studyblog", guestCSRF, guestCSRF, "session-a", 403},
		{"foreign_origin", "POST", "http://attacker.example", "studyblog", guestCSRF, guestCSRF, "", 403},
		{"missing_origin", "POST", "", "studyblog", guestCSRF, guestCSRF, "", 403},
		{"missing_requested_with", "POST", h.config.Origin, "", guestCSRF, guestCSRF, "", 403},
		{"read_without_token", "GET", "", "", "", "", "", 204},
		{"head_without_token", "HEAD", "", "", "", "", "", 204},
		{"options_without_token", "OPTIONS", "", "", "", "", "", 204},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := httptest.NewRequest(tc.method, "/api/posts", nil)
			r.Header.Set("Origin", tc.origin)
			r.Header.Set("X-Requested-With", tc.requestedWith)
			r.Header.Set("X-CSRF-Token", tc.header)
			if tc.cookie != "" {
				r.AddCookie(&http.Cookie{Name: "csrf_token", Value: tc.cookie})
			}
			if tc.session != "" {
				r.AddCookie(&http.Cookie{Name: "access_token", Value: tc.session})
			}
			w := httptest.NewRecorder()
			h.security(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})).ServeHTTP(w, r)
			if w.Code != tc.want {
				t.Errorf("status = %d, want %d", w.Code, tc.want)
			}
		})
	}
}

// Все маршруты изменения данных проверяют CSRF раньше сессии и обращения к БД.
func TestRoutesRejectMissingCSRF(t *testing.T) {
	t.Parallel()
	for _, route := range []struct{ method, path string }{
		{"POST", "/api/auth/login"}, {"POST", "/api/auth/register"},
		{"POST", "/api/auth/logout"}, {"POST", "/api/posts"},
		{"POST", "/api/posts/1/comments"}, {"DELETE", "/api/posts/1"},
		{"PATCH", "/api/profile"}, {"PATCH", "/api/admin/users/1/ban"},
	} {
		t.Run(route.path, func(t *testing.T) {
			t.Parallel()
			r := httptest.NewRequest(route.method, route.path, nil)
			r.Header.Set("Origin", "http://localhost:3000")
			r.Header.Set("X-Requested-With", "studyblog")
			w := httptest.NewRecorder()
			New(nil, config.Config{Origin: "http://localhost:3000"}).ServeHTTP(w, r)
			if w.Code != http.StatusForbidden {
				t.Errorf("status = %d, want 403", w.Code)
			}
		})
	}
}

func TestSessionAllowsNewLoginWithStaleCookie(t *testing.T) {
	t.Parallel()
	for _, path := range []string{"/api/auth/login", "/api/auth/register", "/api/posts"} {
		t.Run(path, func(t *testing.T) {
			t.Parallel()
			r := httptest.NewRequest("POST", path, nil)
			r.AddCookie(&http.Cookie{Name: "access_token", Value: "expired"})
			w := httptest.NewRecorder()
			h := &Handler{}
			h.session(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})).ServeHTTP(w, r)
			want := http.StatusNoContent
			if path == "/api/posts" {
				want = http.StatusUnauthorized
			}
			if w.Code != want {
				t.Errorf("status = %d, want %d", w.Code, want)
			}
		})
	}
}
