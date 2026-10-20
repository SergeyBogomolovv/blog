// Шифрование персональных данных и закрытый поисковый индекс логинов.
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
	"strings"
)

// DataCipher защищает значения AES-256-GCM и индексирует логины отдельным ключом.
type DataCipher struct {
	aead     cipher.AEAD
	indexKey []byte
}

// NewDataCipher принимает два независимых 256-битных ключа из конфигурации.
func NewDataCipher(encryptionKey, indexKey []byte) (*DataCipher, error) {
	if len(encryptionKey) != 32 || len(indexKey) != 32 {
		return nil, errors.New("ключи защиты данных должны содержать по 32 байта")
	}
	if hmac.Equal(encryptionKey, indexKey) {
		return nil, errors.New("ключ шифрования и ключ индекса должны различаться")
	}
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &DataCipher{aead: aead, indexKey: append([]byte(nil), indexKey...)}, nil
}

// Encrypt шифрует значение со случайным nonce и привязкой к назначению поля.
func (c *DataCipher) Encrypt(value, field string) (string, error) {
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	sealed := c.aead.Seal(nonce, nonce, []byte(value), []byte(field))
	return "v1:" + base64.RawStdEncoding.EncodeToString(sealed), nil
}

// Decrypt проверяет целостность значения и запрещает подмену между полями.
func (c *DataCipher) Decrypt(value, field string) (string, error) {
	const message = "не удалось расшифровать персональные данные"
	if !strings.HasPrefix(value, "v1:") {
		return "", errors.New(message)
	}
	data, err := base64.RawStdEncoding.DecodeString(strings.TrimPrefix(value, "v1:"))
	if err != nil || len(data) < c.aead.NonceSize()+c.aead.Overhead() {
		return "", errors.New(message)
	}
	nonce := data[:c.aead.NonceSize()]
	plain, err := c.aead.Open(nil, nonce, data[c.aead.NonceSize():], []byte(field))
	if err != nil {
		return "", errors.New(message)
	}
	return string(plain), nil
}

// LoginIndex строит HMAC нормализованного логина для поиска и проверки уникальности.
func (c *DataCipher) LoginIndex(login string) string {
	mac := hmac.New(sha256.New, c.indexKey)
	mac.Write([]byte("cleanhome:login:v1:" + strings.ToLower(strings.TrimSpace(login))))
	return hex.EncodeToString(mac.Sum(nil))
}
