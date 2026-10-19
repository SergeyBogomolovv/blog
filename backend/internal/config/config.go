// Package config загружает и проверяет настройки из корневого файла .env.
package config

import (
	"encoding/hex"
	"fmt"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN           string
	Address       string
	Origin        string
	JWTKey        []byte
	EncryptionKey []byte
	IndexKey      []byte
	Secure        bool
}

// Load читает настройки ограниченной роли сервера или административной утилиты.
func Load(admin bool) (Config, error) {
	for _, path := range []string{".env", "../.env"} {
		if _, err := os.Stat(path); err == nil {
			if err := godotenv.Load(path); err != nil {
				return Config{}, fmt.Errorf("не удалось прочитать .env")
			}

			break
		}
	}

	c := Config{
		Address: os.Getenv("SERVER_ADDR"),
		Origin:  os.Getenv("APP_ORIGIN"),
		Secure:  os.Getenv("COOKIE_SECURE") == "true",
	}
	for name, target := range map[string]*[]byte{
		"JWT_SECRET":      &c.JWTKey,
		"ENCRYPTION_KEY":  &c.EncryptionKey,
		"LOGIN_INDEX_KEY": &c.IndexKey,
	} {
		key, err := hex.DecodeString(os.Getenv(name))
		if err != nil || len(key) != 32 {
			return c, fmt.Errorf("%s должен содержать 32 байта в hex; выполните scripts/generate-secrets.sh", name)
		}

		*target = key
	}

	role, password := "user", os.Getenv("DB_USER_PASSWORD")
	if admin {
		role, password = "admin", os.Getenv("DB_ADMIN_PASSWORD")
	}

	for _, key := range []string{"DB_HOST", "DB_PORT", "POSTGRES_DB", "SERVER_ADDR", "APP_ORIGIN"} {
		if os.Getenv(key) == "" {
			return c, fmt.Errorf("не задано значение %s", key)
		}
	}

	if password == "" {
		return c, fmt.Errorf("не задан пароль базы данных")
	}

	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(role, password),
		Host:   os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		Path:   "/" + os.Getenv("POSTGRES_DB"),
	}
	q := u.Query()
	q.Set("sslmode", "disable")
	u.RawQuery = q.Encode()
	c.DSN = u.String()
	return c, nil
}
