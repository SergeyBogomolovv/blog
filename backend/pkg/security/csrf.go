// Подпись CSRF-токена совместима с серверным модулем Next.js lib/csrf.ts.
package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// ValidCSRF проверяет HMAC и привязку к access-cookie, в том числе для гостя.
func ValidCSRF(key []byte, token, session string) bool {
	if len(token) != 129 {
		return false
	}
	nonce, signature, ok := strings.Cut(token, ".")
	if !ok || len(nonce) != 64 || len(signature) != 64 {
		return false
	}
	if _, err := hex.DecodeString(nonce); err != nil {
		return false
	}
	provided, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, key)
	fmt.Fprintf(mac, "studyblog:csrf:v1!%d!%s!%s", len(session), session, nonce)
	return hmac.Equal(provided, mac.Sum(nil))
}
