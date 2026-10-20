// Единые серверные правила регистрации, поиска и заявки.
package service

import (
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"cleanhome/internal/model"
)

const (
	ConsentVersion    = "2026-10-20"
	maxTextLength     = 20
	maxAddressLength  = 100
	minPasswordLength = 8
	dateLayout        = "2006-01-02"
	domainLabel = `[A-Za-z0-9]([A-Za-z0-9\-]*[A-Za-z0-9])?`
)

var (
	moscowTime = time.FixedZone("Europe/Moscow", 3*60*60)

	surnamePattern    = regexp.MustCompile(`^[А-ЯЁ][а-яё]+$`)
	namePattern       = regexp.MustCompile(`^[А-ЯЁ][а-яё]+(-[А-ЯЁ][а-яё]+)?$`)
	patronymicPattern = regexp.MustCompile(`^[А-ЯЁ][а-яё]+([ -](оглы|кызы|уулу|гызы))?$`)
	loginPattern      = regexp.MustCompile(`^[A-Za-z]+$`)
	phonePattern      = regexp.MustCompile(`^8\([0-9]{3}\)[0-9]{3}-[0-9]{2}-[0-9]{2}$`)
	searchPattern     = regexp.MustCompile(`^[А-Яа-яЁёA-Za-z0-9 -]+$`)

	emailNamePattern   = regexp.MustCompile(`^[A-Za-z0-9]([A-Za-z0-9._%+\-]*[A-Za-z0-9_%+\-])?$`)
	emailDomainPattern = regexp.MustCompile(`^` + domainLabel + `(\.` + domainLabel + `)+$`)

	uppercasePattern = regexp.MustCompile(`[A-Z]`)
	lowercasePattern = regexp.MustCompile(`[a-z]`)
	digitPattern     = regexp.MustCompile(`[0-9]`)
	specialPattern   = regexp.MustCompile(`[$#@!]`)
)

var cleaningServices = []model.CleaningService{
	{
		ID:   "maintenance",
		Name: "Поддерживающая уборка",
		Description: "Регулярная забота о доме: убираем пыль, моем полы " +
			"и приводим в порядок кухню и ванную.",
	},
	{
		ID:   "deep",
		Name: "Генеральная уборка",
		Description: "Тщательная уборка всего дома, включая труднодоступные места, " +
			"поверхности и накопившиеся загрязнения.",
	},
	{
		ID:   "renovation",
		Name: "Уборка после ремонта",
		Description: "Удаляем строительную пыль и следы ремонта. " +
			"Подготавливаем чистое пространство к новой жизни.",
	},
}

type Fields map[string]string

func (fields Fields) Error() string {
	return "Проверьте заполнение полей."
}

type Registration struct {
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	Patronymic string `json:"patronymic"`
	Login      string `json:"login"`
	Password   string `json:"password"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Consent    bool   `json:"consent"`
}

func (input *Registration) Validate() Fields {
	input.Surname = strings.TrimSpace(input.Surname)
	input.Name = strings.TrimSpace(input.Name)
	input.Patronymic = strings.TrimSpace(input.Patronymic)
	input.Login = strings.TrimSpace(input.Login)
	input.Phone = strings.TrimSpace(input.Phone)
	input.Email = strings.TrimSpace(input.Email)

	fields := Fields{}
	if !validText(input.Surname, surnamePattern) {
		fields["surname"] = "Фамилия: русские буквы, первая заглавная, от 2 до 20 символов."
	}
	if !validText(input.Name, namePattern) {
		fields["name"] = "Имя: от 2 до 20 русских букв. Допустимо Анна-Мария."
	}
	if !validText(input.Patronymic, patronymicPattern) {
		fields["patronymic"] = "Отчество с заглавной буквы, до 20 символов. Допустимо Рашид оглы."
	}
	if !validText(input.Login, loginPattern) {
		fields["login"] = "Логин: только латинские буквы, от 1 до 20 символов."
	}
	if !validText(input.Phone, phonePattern) {
		fields["phone"] = "Введите телефон в формате 8(999)123-45-67."
	}
	if !validEmail(input.Email) {
		fields["email"] = "Email до 20 символов: name@mail.ru. Проверьте точки в имени и домен."
	}
	if !validPassword(input.Password) {
		fields["password"] = "От 8 символов: заглавная и строчная латинские буквы, цифра и символ $#@!."
	}
	if !input.Consent {
		fields["consent"] = "Для регистрации подтвердите согласие на обработку данных."
	}

	return fields
}

func validText(value string, pattern *regexp.Regexp) bool {
	length := utf8.RuneCountInString(value)
	return length >= 1 && length <= maxTextLength && pattern.MatchString(value)
}

func validEmail(email string) bool {
	if length := utf8.RuneCountInString(email); length < 1 || length > maxTextLength {
		return false
	}

	name, domain, hasSeparator := strings.Cut(email, "@")
	return hasSeparator &&
		emailNamePattern.MatchString(name) &&
		emailDomainPattern.MatchString(domain) &&
		!strings.Contains(name, "..")
}

func validPassword(password string) bool {
	return utf8.RuneCountInString(password) >= minPasswordLength &&
		uppercasePattern.MatchString(password) &&
		lowercasePattern.MatchString(password) &&
		digitPattern.MatchString(password) &&
		specialPattern.MatchString(password)
}

func Search(query string) ([]model.CleaningService, error) {
	query = strings.TrimSpace(query)
	if query != "" && !validText(query, searchPattern) {
		return nil, Fields{
			"q": "От 1 до 20 символов: буквы, цифры, пробелы и дефисы.",
		}
	}

	query = strings.ToLower(query)
	matches := make([]model.CleaningService, 0, len(cleaningServices))
	for _, cleaningService := range cleaningServices {
		text := strings.ToLower(cleaningService.Name + " " + cleaningService.Description)
		if strings.Contains(text, query) {
			matches = append(matches, cleaningService)
		}
	}

	return matches, nil
}

type RequestInput struct {
	Service string `json:"service"`
	Address string `json:"address"`
	Date    string `json:"date"`
}

func (input *RequestInput) Validate(now time.Time) Fields {
	input.Service = strings.TrimSpace(input.Service)
	input.Address = strings.TrimSpace(input.Address)
	input.Date = strings.TrimSpace(input.Date)

	fields := Fields{}
	if !serviceExists(input.Service) {
		fields["service"] = "Выберите услугу из списка."
	}
	if length := utf8.RuneCountInString(input.Address); length < 1 || length > maxAddressLength {
		fields["address"] = "Адрес должен содержать от 1 до 100 символов."
	}
	if !validCleaningDate(input.Date, now) {
		fields["date"] = "Выберите существующую дату не раньше сегодняшнего дня по Москве."
	}

	return fields
}

func serviceExists(serviceID string) bool {
	for _, cleaningService := range cleaningServices {
		if cleaningService.ID == serviceID {
			return true
		}
	}

	return false
}

func validCleaningDate(value string, now time.Time) bool {
	date, err := time.Parse(dateLayout, value)
	if err != nil || date.Format(dateLayout) != value {
		return false
	}

	return value >= MoscowToday(now)
}

func MoscowToday(now time.Time) string {
	return now.In(moscowTime).Format(dateLayout)
}
