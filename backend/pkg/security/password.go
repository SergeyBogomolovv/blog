// Хеширование паролей с индивидуальной солью и фиксированными параметрами scrypt.
package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"strings"

	"golang.org/x/crypto/scrypt"
)

const (
	passwordSaltSize = 16
	passwordKeySize  = 32
	scryptCost        = 32768
	scryptBlockSize   = 8
	scryptParallelism = 1
	passwordHashType  = "scrypt$32768$8$1"
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, passwordSaltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	key, err := derivePasswordKey(password, salt)
	if err != nil {
		return "", err
	}

	return strings.Join([]string{
		passwordHashType,
		hex.EncodeToString(salt),
		hex.EncodeToString(key),
	}, "$"), nil
}

func CheckPassword(hash, password string) bool {
	parts := strings.Split(hash, "$")
	if len(parts) != 6 || strings.Join(parts[:4], "$") != passwordHashType {
		return false
	}

	salt, err := hex.DecodeString(parts[4])
	if err != nil || len(salt) != passwordSaltSize {
		return false
	}

	expected, err := hex.DecodeString(parts[5])
	if err != nil || len(expected) != passwordKeySize {
		return false
	}

	actual, err := derivePasswordKey(password, salt)
	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare(actual, expected) == 1
}

func derivePasswordKey(password string, salt []byte) ([]byte, error) {
	return scrypt.Key(
		[]byte(password),
		salt,
		scryptCost,
		scryptBlockSize,
		scryptParallelism,
		passwordKeySize,
	)
}
