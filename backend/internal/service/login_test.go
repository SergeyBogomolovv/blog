// Проверки отказа при неверном формате входа до обращения к хранилищу.
package service

import (
	"context"
	"errors"
	"testing"
)

// TestInvalidLogin проверяет, что неверные учётные данные не доходят до БД.
func TestInvalidLogin(t *testing.T) {
	s := &Service{}
	for _, credentials := range [][2]string{
		{"", "Password1!"},
		{"123", "Password1!"},
		{"Иван", "Password1!"},
		{"abcdefghijklmnopqrstu", "Password1!"},
		{"Ivan", ""},
		{"Ivan", "password"},
	} {
		if _, err := s.Login(context.Background(), credentials[0], credentials[1]); !errors.Is(err, ErrCredentials) {
			t.Fatalf("ожидался отказ для неверного формата: %v", err)
		}
	}
}
