// Проверки конфиденциальности, целостности и поискового индекса персональных данных.
package security

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"
)

// TestDataEncryption проверяет случайность шифротекста и восстановление русских строк.
func TestDataEncryption(t *testing.T) {
	c, err := NewDataCipher(bytes.Repeat([]byte{1}, 32), bytes.Repeat([]byte{2}, 32))
	if err != nil {
		t.Fatal(err)
	}
	const value = "Анна-Мария, Москва, дом 7"
	first, err := c.Encrypt(value, "users.name")
	if err != nil {
		t.Fatal(err)
	}
	second, err := c.Encrypt(value, "users.name")
	if err != nil {
		t.Fatal(err)
	}
	if first == second || strings.Contains(first, value) {
		t.Fatal("шифрование должно скрывать текст и использовать новый nonce")
	}
	plain, err := c.Decrypt(first, "users.name")
	if err != nil || plain != value {
		t.Fatalf("данные не восстановлены: %v", err)
	}

	raw, err := base64.RawStdEncoding.DecodeString(strings.TrimPrefix(first, "v1:"))
	if err != nil {
		t.Fatal(err)
	}
	raw[len(raw)-1] ^= 1
	tampered := "v1:" + base64.RawStdEncoding.EncodeToString(raw)
	for _, ciphertext := range []string{tampered, "v1:AA", "открытый текст", "v1:???"} {
		if _, err := c.Decrypt(ciphertext, "users.name"); err == nil {
			t.Fatal("повреждённый или незашифрованный текст принят")
		}
	}
	if _, err := c.Decrypt(first, "users.email"); err == nil {
		t.Fatal("подмена между полями не обнаружена")
	}
	other, err := NewDataCipher(bytes.Repeat([]byte{3}, 32), bytes.Repeat([]byte{2}, 32))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := other.Decrypt(first, "users.name"); err == nil {
		t.Fatal("данные расшифрованы чужим ключом")
	}
}

// TestLoginIndex проверяет нормализацию логина и зависимость индекса от секретного ключа.
func TestLoginIndex(t *testing.T) {
	c, err := NewDataCipher(bytes.Repeat([]byte{1}, 32), bytes.Repeat([]byte{2}, 32))
	if err != nil {
		t.Fatal(err)
	}
	index := c.LoginIndex("Ivan")
	if index != c.LoginIndex(" ivAN ") || len(index) != 64 {
		t.Fatal("индекс должен учитывать нормализацию логина")
	}
	if index == c.LoginIndex("Anna") {
		t.Fatal("разные логины получили одинаковый индекс")
	}
	other, err := NewDataCipher(bytes.Repeat([]byte{1}, 32), bytes.Repeat([]byte{3}, 32))
	if err != nil {
		t.Fatal(err)
	}
	if index == other.LoginIndex("Ivan") {
		t.Fatal("индекс не зависит от ключа")
	}
	for _, keys := range [][2][]byte{
		{nil, bytes.Repeat([]byte{2}, 32)},
		{bytes.Repeat([]byte{1}, 32), nil},
		{bytes.Repeat([]byte{1}, 32), bytes.Repeat([]byte{1}, 32)},
	} {
		if _, err := NewDataCipher(keys[0], keys[1]); err == nil {
			t.Fatal("недопустимые ключи приняты")
		}
	}
}
