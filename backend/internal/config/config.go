// Package config читает настройки локального сервиса из общего .env.
package config

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

// Config параметры подключения, адрес сайта и ключи защиты.
type Config struct {
	DSN        string
	Address    string
	Origin     string
	SessionKey []byte
	DataKey    []byte
	IndexKey   []byte
}

// Load загружает и проверяет настройки выбранной роли приложения.
func Load(admin bool) (Config, error) {
	if err := loadEnvFile(); err != nil {
		return Config{}, err
	}

	key, err := hex.DecodeString(os.Getenv("SESSION_SECRET"))
	if err != nil || len(key) != 32 {
		return Config{}, errors.New(
			"SESSION_SECRET должен содержать 32 байта в hex; выполните make secrets",
		)
	}
	dataKey, err := readKey("DATA_ENCRYPTION_KEY")
	if err != nil {
		return Config{}, err
	}
	indexKey, err := readKey("LOGIN_INDEX_KEY")
	if err != nil {
		return Config{}, err
	}

	required := []string{
		"DB_HOST",
		"DB_PORT",
		"POSTGRES_DB",
		"SERVER_ADDR",
		"APP_ORIGIN",
	}
	for _, name := range required {
		if os.Getenv(name) == "" {
			return Config{}, fmt.Errorf("не задано значение %s", name)
		}
	}

	dsn, err := databaseDSN(admin)
	if err != nil {
		return Config{}, err
	}

	return Config{
		DSN:        dsn,
		Address:    os.Getenv("SERVER_ADDR"),
		Origin:     os.Getenv("APP_ORIGIN"),
		SessionKey: key,
		DataKey:    dataKey,
		IndexKey:   indexKey,
	}, nil
}

// readKey читает 256-битный ключ, не включая его значение в сообщения об ошибках.
func readKey(name string) ([]byte, error) {
	key, err := hex.DecodeString(os.Getenv(name))
	if err != nil || len(key) != 32 {
		return nil, fmt.Errorf("%s должен содержать 32 байта в hex", name)
	}
	return key, nil
}

// loadEnvFile читает локальную конфигурацию, сохраняя приоритет окружения.
func loadEnvFile() error {
	for _, path := range []string{".env", "../.env"} {
		if _, err := os.Stat(path); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}

			return fmt.Errorf("не удалось открыть .env: %w", err)
		}
		if err := godotenv.Load(path); err != nil {
			return fmt.Errorf("не удалось прочитать .env: %w", err)
		}

		return nil
	}

	// При запуске без файла настройки могут быть заданы в окружении процесса.
	return nil
}

// databaseDSN собирает строку подключения с паролем из окружения.
func databaseDSN(admin bool) (string, error) {
	role := "user"
	password := os.Getenv("DB_USER_PASSWORD")
	if admin {
		role = "admin"
		password = os.Getenv("DB_ADMIN_PASSWORD")
	}
	if password == "" {
		return "", errors.New("не задан пароль базы данных")
	}

	address := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(role, password),
		Host:   os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		Path:   "/" + os.Getenv("POSTGRES_DB"),
	}
	parameters := address.Query()
	parameters.Set("sslmode", "disable")
	address.RawQuery = parameters.Encode()

	return address.String(), nil
}
