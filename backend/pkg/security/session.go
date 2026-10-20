// Случайные идентификаторы сессий подписываются перед записью в cookie.
package security

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

const sessionIDSize = 32

// SessionID создаёт случайный идентификатор сессии с подписью HMAC.
func SessionID(key []byte) (string, error) {
	id := make([]byte, sessionIDSize)
	if _, err := rand.Read(id); err != nil {
		return "", err
	}

	value := hex.EncodeToString(id)
	signature := hex.EncodeToString(sessionMAC(key, value))
	return value + "." + signature, nil
}

// ValidSession проверяет формат и подпись токена сессии.
func ValidSession(key []byte, value string) bool {
	id, signature, hasSeparator := strings.Cut(value, ".")
	if !hasSeparator || len(id) != sessionIDSize*2 || len(signature) != sha256.Size*2 {
		return false
	}
	if _, err := hex.DecodeString(id); err != nil {
		return false
	}

	provided, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	return hmac.Equal(provided, sessionMAC(key, id))
}

// SessionHash преобразует токен в безопасный для хранения SHA-256.
func SessionHash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

// sessionMAC подписывает идентификатор сессии серверным ключом.
func sessionMAC(key []byte, id string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte("cleaning:session:v1!" + id))
	return mac.Sum(nil)
}
