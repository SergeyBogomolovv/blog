// Package security реализует шифрование полей, индексы логинов и JWT.
package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Vault struct {
	aead     cipher.AEAD
	indexKey []byte
}

type Claims struct {
	Version int `json:"ver"`
	jwt.RegisteredClaims
}

// NewVault создаёт шифр AES-256-GCM с отдельным ключом индекса логинов.
func NewVault(key, indexKey []byte) (*Vault, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	return &Vault{aead, indexKey}, err
}

// Encrypt шифрует поле с новым случайным одноразовым значением.
func (v *Vault) Encrypt(value string) (string, error) {
	nonce := make([]byte, v.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(v.aead.Seal(nonce, nonce, []byte(value), nil)), nil
}

// Decrypt проверяет целостность шифротекста перед расшифровкой.
func (v *Vault) Decrypt(value string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil || len(data) < v.aead.NonceSize() {
		return "", errors.New("invalid ciphertext")
	}

	n := v.aead.NonceSize()
	plain, err := v.aead.Open(nil, data[:n], data[n:], nil)
	return string(plain), err
}

// Index создаёт индекс поиска без хранения логина в открытом виде.
func (v *Vault) Index(login string) string {
	mac := hmac.New(sha256.New, v.indexKey)
	mac.Write([]byte(login))
	return hex.EncodeToString(mac.Sum(nil))
}

// Issue подписывает двухчасовой токен с идентификатором и версией сессии.
func Issue(key []byte, id int64, version int) (string, error) {
	now := time.Now()
	c := Claims{
		Version:   version,
		Subject:   strconv.FormatInt(id, 10),
		Issuer:    "studyblog",
		Audience:  jwt.ClaimStrings{"studyblog"},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(2 * time.Hour)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(key)
}

// Parse принимает только действующие HS256-токены этого приложения.
func Parse(key []byte, raw string) (*Claims, error) {
	c := &Claims{}
	_, err := jwt.ParseWithClaims(
		raw,
		c,
		func(t *jwt.Token) (any, error) {
			return key, nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer("studyblog"),
		jwt.WithAudience("studyblog"),
		jwt.WithExpirationRequired(),
	)
	return c, err
}
