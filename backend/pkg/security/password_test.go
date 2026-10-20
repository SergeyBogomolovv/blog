package security_test

import (
	"bytes"
	"strings"
	"testing"

	"cleanhome/pkg/security"
)

func TestPasswords(t *testing.T) {
	password := "Aa1!" + strings.Repeat("я", 100)
	firstHash, err := security.HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}

	secondHash, err := security.HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	if firstHash == secondHash {
		t.Fatal("the same password must receive a different salt each time")
	}

	cases := []struct {
		name     string
		hash     string
		password string
		valid    bool
	}{
		{name: "correct password", hash: firstHash, password: password, valid: true},
		{name: "wrong password", hash: firstHash, password: "wrong"},
		{name: "trailing space matters", hash: firstHash, password: password + " "},
		{name: "malformed hash", hash: "invalid", password: password},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if valid := security.CheckPassword(tc.hash, tc.password); valid != tc.valid {
				t.Fatalf("CheckPassword() = %v, want %v", valid, tc.valid)
			}
		})
	}
}

func TestSignedSession(t *testing.T) {
	key := bytes.Repeat([]byte{1}, 32)
	token, err := security.SessionID(key)
	if err != nil {
		t.Fatal(err)
	}
	if !security.ValidSession(key, token) {
		t.Fatal("a freshly signed session must be valid")
	}

	otherKey := bytes.Repeat([]byte{2}, 32)
	if security.ValidSession(otherKey, token) {
		t.Fatal("a different key must not validate the session")
	}
	if security.ValidSession(key, token+"x") {
		t.Fatal("a modified token must not validate")
	}

	otherToken, err := security.SessionID(key)
	if err != nil {
		t.Fatal(err)
	}
	if otherToken == token {
		t.Fatal("session identifiers must not be reused")
	}
	if security.SessionHash(otherToken) == security.SessionHash(token) {
		t.Fatal("different identifiers must produce different stored hashes")
	}
}
