// Браузер повторяет правила Go; окончательное решение принимает сервер.
import type {
  CleaningService,
  FieldErrors,
  RegistrationValues,
  RequestValues,
} from "./api.ts";

export const consentVersion = "2026-10-20";
export const passwordHint =
  "От 8 символов: A–Z, a–z, цифра и один из символов $#@!.";

const textPatterns = {
  surname: /^[А-ЯЁ][а-яё]+$/,
  name: /^[А-ЯЁ][а-яё]+(-[А-ЯЁ][а-яё]+)?$/,
  patronymic: /^[А-ЯЁ][а-яё]+([ -](оглы|кызы|уулу|гызы))?$/,
  login: /^[A-Za-z]+$/,
  phone: /^8\([0-9]{3}\)[0-9]{3}-[0-9]{2}-[0-9]{2}$/,
};

const emailNamePattern = /^[A-Za-z0-9]([A-Za-z0-9._%+\-]*[A-Za-z0-9_%+\-])?$/;
const domainLabel = "[A-Za-z0-9]([A-Za-z0-9\\-]*[A-Za-z0-9])?";
const emailDomainPattern = new RegExp(`^${domainLabel}(\\.${domainLabel})+$`);
const searchPattern = /^[А-Яа-яЁёA-Za-z0-9 -]+$/;
const calendarDatePattern = /^\d{4}-\d{2}-\d{2}$/;

/** Считает символы Unicode независимо от числа кодовых единиц. */
function characterCount(value: string): number {
  return Array.from(value).length;
}

/** Проверяет длину строки и соответствие допустимому формату. */
function validText(value: string, pattern: RegExp): boolean {
  const length = characterCount(value);

  return length >= 1 && length <= 20 && pattern.test(value);
}

/** Проверяет минимальную длину и обязательные символы пароля. */
export function validPassword(password: string): boolean {
  return (
    characterCount(password) >= 8 &&
    /[A-Z]/.test(password) &&
    /[a-z]/.test(password) &&
    /[0-9]/.test(password) &&
    /[$#@!]/.test(password)
  );
}

/** Проверяет длину, домен и допустимые точки в адресе почты. */
export function validEmail(email: string): boolean {
  if (characterCount(email) > 20) {
    return false;
  }

  const parts = email.split("@");

  if (parts.length !== 2) {
    return false;
  }

  const [name, domain] = parts;

  return (
    emailNamePattern.test(name) &&
    emailDomainPattern.test(domain) &&
    !name.includes("..")
  );
}

/** Проверяет длину и разрешённые символы поискового запроса. */
export function validSearch(query: string): boolean {
  return validText(query, searchPattern);
}

/** Проверяет данные входа до передачи на сервер. */
export function validateLogin(input: {
  login: string;
  password: string;
}): FieldErrors {
  const fields: FieldErrors = {};

  if (!validText(input.login.trim(), textPatterns.login)) {
    fields.login = "Логин: только латинские буквы, от 1 до 20 символов.";
  }

  if (!validPassword(input.password)) {
    fields.password = "Логин или пароль некорректен";
  }

  return fields;
}

/** Возвращает ошибки регистрационных полей и согласия. */
export function validateRegistration(input: RegistrationValues): FieldErrors {
  const fields: FieldErrors = {};

  for (const field of [
    "surname",
    "name",
    "patronymic",
    "login",
    "phone",
  ] as const) {
    if (!validText(input[field].trim(), textPatterns[field])) {
      fields[field] = "Проверьте формат и длину поля (до 20 символов).";
    }
  }

  if (!validEmail(input.email.trim())) {
    fields.email = "Введите email до 20 символов, например name@mail.ru.";
  }

  if (!validPassword(input.password)) {
    fields.password = passwordHint;
  }

  if (!input.consent) {
    fields.consent = "Подтвердите согласие на обработку данных.";
  }

  return fields;
}

/** Отклоняет несуществующие и прошедшие календарные даты. */
function validCleaningDate(value: string, today: string): boolean {
  if (!calendarDatePattern.test(value)) {
    return false;
  }

  const parsed = new Date(`${value}T00:00:00Z`);

  if (Number.isNaN(parsed.getTime())) {
    return false;
  }

  return parsed.toISOString().slice(0, 10) === value && value >= today;
}

/** Проверяет услугу, адрес и дату перед отправкой заявки. */
export function validateRequest(
  input: RequestValues,
  services: CleaningService[],
  today: string,
): FieldErrors {
  const fields: FieldErrors = {};
  const hasService = services.some((service) => service.id === input.service);
  const addressLength = characterCount(input.address.trim());

  if (!hasService) {
    fields.service = "Выберите услугу.";
  }

  if (addressLength < 1 || addressLength > 100) {
    fields.address = "Адрес должен содержать от 1 до 100 символов.";
  }

  if (!validCleaningDate(input.date, today)) {
    fields.date =
      "Выберите существующую дату не раньше сегодняшнего дня по Москве.";
  }

  return fields;
}
