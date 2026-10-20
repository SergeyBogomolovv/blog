// Подпись CSRF-токена совместима с серверным модулем Next.js lib/csrf.ts.
package security

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

const csrfNonceSize = 32

// ValidCSRF проверяет подпись и привязку к cookie сессии, в том числе для гостя.
func ValidCSRF(key []byte, token, session string) bool {
	const tokenLength = csrfNonceSize*2 + 1 + sha256.Size*2
	if len(token) != tokenLength {
		return false
	}

	nonce, signature, hasSeparator := strings.Cut(token, ".")
	if !hasSeparator || len(nonce) != csrfNonceSize*2 || len(signature) != sha256.Size*2 {
		return false
	}
	if _, err := hex.DecodeString(nonce); err != nil {
		return false
	}

	provided, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	return hmac.Equal(provided, csrfMAC(key, nonce, session))
}

// NewCSRF подписывает случайный токен, привязанный к текущей сессии.
func NewCSRF(key []byte, session string) (string, error) {
	nonce := make([]byte, csrfNonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	value := hex.EncodeToString(nonce)
	signature := hex.EncodeToString(csrfMAC(key, value, session))
	return value + "." + signature, nil
}

// csrfMAC вычисляет подпись CSRF с отдельным контекстом назначения.
func csrfMAC(key []byte, nonce, session string) []byte {
	mac := hmac.New(sha256.New, key)
	fmt.Fprintf(mac, "cleanhome:csrf:v1!%d!%s!%s", len(session), session, nonce)
	return mac.Sum(nil)
}
