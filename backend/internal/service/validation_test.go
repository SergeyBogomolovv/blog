package service_test

import (
	"strings"
	"testing"
	"time"

	"cleanhome/internal/service"
)

func validRegistration() service.Registration {
	return service.Registration{
		Surname:    "Иванов",
		Name:       "Анна-Мария",
		Patronymic: "Рашид оглы",
		Login:      "Ivan",
		Password:   "Password1!",
		Phone:      "8(999)123-45-67",
		Email:      "a@mail.ru",
		Consent:    true,
	}
}

func TestRegistration(t *testing.T) {
	tests := []struct {
		name   string
		field  string
		change func(*service.Registration)
	}{
		{
			name:  "lowercase surname",
			field: "surname",
			change: func(input *service.Registration) {
				input.Surname = "иванов"
			},
		},
		{
			name:  "uppercase name",
			field: "name",
			change: func(input *service.Registration) {
				input.Name = "АННА"
			},
		},
		{
			name:  "name 21",
			field: "name",
			change: func(input *service.Registration) {
				input.Name = "А" + strings.Repeat("а", 20)
			},
		},
		{
			name:  "login digit",
			field: "login",
			change: func(input *service.Registration) {
				input.Login = "Ivan1"
			},
		},
		{
			name:  "login 21",
			field: "login",
			change: func(input *service.Registration) {
				input.Login = strings.Repeat("a", 21)
			},
		},
		{
			name:  "phone",
			field: "phone",
			change: func(input *service.Registration) {
				input.Phone = "+79991234567"
			},
		},
		{
			name:  "email leading dot",
			field: "email",
			change: func(input *service.Registration) {
				input.Email = ".a@mail.ru"
			},
		},
		{
			name:  "email trailing dot",
			field: "email",
			change: func(input *service.Registration) {
				input.Email = "a.@mail.ru"
			},
		},
		{
			name:  "email double dot",
			field: "email",
			change: func(input *service.Registration) {
				input.Email = "a..b@mail.ru"
			},
		},
		{
			name:  "email missing domain",
			field: "email",
			change: func(input *service.Registration) {
				input.Email = "a@mail"
			},
		},
		{
			name:  "email 21",
			field: "email",
			change: func(input *service.Registration) {
				input.Email = strings.Repeat("a", 16) + "@b.ru"
			},
		},
		{
			name:  "short password",
			field: "password",
			change: func(input *service.Registration) {
				input.Password = "Aa1!"
			},
		},
		{
			name:  "missing symbol",
			field: "password",
			change: func(input *service.Registration) {
				input.Password = "Password1"
			},
		},
		{
			name:  "no consent",
			field: "consent",
			change: func(input *service.Registration) {
				input.Consent = false
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := validRegistration()
			tc.change(&input)

			if fields := input.Validate(); fields[tc.field] == "" {
				t.Fatalf("expected an error for %s, got %v", tc.field, fields)
			}
		})
	}

	t.Run("valid boundaries and trimming", func(t *testing.T) {
		input := validRegistration()
		input.Name = " А" + strings.Repeat("а", 19) + " "
		input.Login = strings.Repeat("a", 20)
		input.Email = strings.Repeat("a", 15) + "@b.ru"
		input.Password = "Aa1!" + strings.Repeat("x", 100)

		if fields := input.Validate(); len(fields) > 0 {
			t.Fatal(fields)
		}
		if strings.HasPrefix(input.Name, " ") {
			t.Fatal("name was not trimmed")
		}
	})

	patronymics := []string{
		"Рашид оглы",
		"Рашид-кызы",
		"Рашид уулу",
		"Рашид-гызы",
		"Иванович",
	}
	for _, patronymic := range patronymics {
		t.Run(patronymic, func(t *testing.T) {
			input := validRegistration()
			input.Patronymic = patronymic

			if fields := input.Validate(); len(fields) > 0 {
				t.Fatal(fields)
			}
		})
	}
}

func TestRequests(t *testing.T) {
	now := time.Date(2026, 10, 20, 22, 0, 0, 0, time.UTC)
	if today := service.MoscowToday(now); today != "2026-10-21" {
		t.Fatalf("MoscowToday() = %s, want 2026-10-21", today)
	}

	tests := []struct {
		name    string
		address string
		date    string
		service string
		field   string
	}{
		{
			name:    "100 valid",
			address: strings.Repeat("я", 100),
			date:    "2026-10-21",
			service: "deep",
		},
		{
			name:    "101 invalid",
			address: strings.Repeat("я", 101),
			date:    "2026-10-21",
			service: "deep",
			field:   "address",
		},
		{
			name:    "empty address",
			address: " ",
			date:    "2026-10-21",
			service: "deep",
			field:   "address",
		},
		{
			name:    "invalid service",
			address: "Москва",
			date:    "2026-10-21",
			service: "invalid",
			field:   "service",
		},
		{
			name:    "yesterday in Moscow",
			address: "Москва",
			date:    "2026-10-20",
			service: "deep",
			field:   "date",
		},
		{
			name:    "nonexistent date",
			address: "Москва",
			date:    "2027-02-29",
			service: "deep",
			field:   "date",
		},
		{
			name:    "valid leap date",
			address: "Москва",
			date:    "2028-02-29",
			service: "deep",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := service.RequestInput{
				Service: tc.service,
				Address: tc.address,
				Date:    tc.date,
			}
			fields := input.Validate(now)

			if tc.field == "" && len(fields) > 0 {
				t.Fatal(fields)
			}
			if tc.field != "" && fields[tc.field] == "" {
				t.Fatalf("expected an error for %s, got %v", tc.field, fields)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	tests := []struct {
		query   string
		count   int
		invalid bool
	}{
		{query: "", count: 3},
		{query: "РЕМОНТА", count: 1},
		{query: "несуществующее"},
		{query: strings.Repeat("я", 20)},
		{query: strings.Repeat("я", 21), invalid: true},
		{query: "<script>", invalid: true},
		{query: ".*", invalid: true},
		{query: "a|b", invalid: true},
		{query: "[abc]", invalid: true},
	}

	for _, tc := range tests {
		t.Run(tc.query, func(t *testing.T) {
			services, err := service.Search(tc.query)
			if (err != nil) != tc.invalid {
				t.Fatalf("Search() error = %v, want invalid = %v", err, tc.invalid)
			}
			if err == nil && len(services) != tc.count {
				t.Fatalf("got %d services, want %d", len(services), tc.count)
			}
		})
	}
}
