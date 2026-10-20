// Проверки клиентских правил регистрации, входа, поиска и заявок.
import assert from "node:assert/strict";
import { test } from "node:test";
import {
  validateLogin,
  validateRegistration,
  validateRequest,
  validEmail,
  validPassword,
  validSearch,
} from "./validation.ts";

// Проверяет ограничения логина и пароля до отправки формы входа.
test("вход отклоняет неверный формат и принимает допустимые границы", () => {
  assert.deepEqual(
    validateLogin({ login: " Ivan ", password: "Password1!" }),
    {},
  );
  assert.deepEqual(
    validateLogin({
      login: "a".repeat(20),
      password: "Aa1!" + "я".repeat(100),
    }),
    {},
  );

  for (const login of ["", " ", "123", "Иван", "Ivan!", "a".repeat(21)]) {
    assert.ok(validateLogin({ login, password: "Password1!" }).login);
  }

  assert.equal(
    validateLogin({ login: "Ivan", password: "" }).password,
    "Логин или пароль некорректен",
  );
});

const valid = {
  surname: "Иванов",
  name: "Анна-Мария",
  patronymic: "Рашид оглы",
  login: "Ivan",
  password: "Password1!",
  phone: "8(999)123-45-67",
  email: "a@mail.ru",
  consent: true,
};

test("регистрация, допустимые ФИО и границы", () => {
  assert.deepEqual(validateRegistration(valid), {});
  assert.deepEqual(
    validateRegistration({
      ...valid,
      name: "А" + "а".repeat(19),
      login: "a".repeat(20),
      email: "a".repeat(15) + "@b.ru",
      password: "Aa1!" + "я".repeat(100),
    }),
    {},
  );

  for (const [field, value] of [
    ["surname", "иванов"],
    ["name", "АННА"],
    ["name", "А" + "а".repeat(20)],
    ["login", "Ivan1"],
    ["login", "a".repeat(21)],
    ["phone", "+79991234567"],
    ["consent", false],
  ]) {
    assert.ok(validateRegistration({ ...valid, [field]: value })[field]);
  }

  for (const patronymic of [
    "Рашид оглы",
    "Рашид-кызы",
    "Рашид уулу",
    "Рашид-гызы",
  ]) {
    assert.deepEqual(validateRegistration({ ...valid, patronymic }), {});
  }
});

test("email", () => {
  for (const value of ["a@mail.ru", "a.b@mail.ru", "a".repeat(15) + "@b.ru"]) {
    assert.equal(validEmail(value), true, value);
  }

  for (const value of [
    ".a@mail.ru",
    "a.@mail.ru",
    "a..b@mail.ru",
    "a@mail",
    "a".repeat(16) + "@b.ru",
  ]) {
    assert.equal(validEmail(value), false, value);
  }
});

test("пароль без ограничения 20", () => {
  assert.equal(validPassword("Aa1!" + "я".repeat(100)), true);

  for (const value of [
    "Aa1!",
    "password1!",
    "PASSWORD1!",
    "Password1",
    "Password!",
  ]) {
    assert.equal(validPassword(value), false);
  }
});

test("поиск подстрок без интерпретации кода", () => {
  for (const q of ["ремонта", "Уборка-1", "я".repeat(20)]) {
    assert.equal(validSearch(q), true);
  }

  for (const q of ["", "я".repeat(21), "<script>", ".*", "a|b", "[abc]"]) {
    assert.equal(validSearch(q), false);
  }
});

test("дата уборки и границы заявки проверяются до отправки", () => {
  const services = [
    { id: "deep", name: "Генеральная уборка", description: "" },
  ];
  const today = "2026-10-21";
  const input = { service: "deep", address: "Москва", date: today };

  assert.deepEqual(validateRequest(input, services, today), {});
  assert.deepEqual(
    validateRequest({ ...input, address: "я".repeat(100) }, services, today),
    {},
  );
  assert.deepEqual(
    validateRequest({ ...input, date: "2028-02-29" }, services, today),
    {},
  );

  const invalidValues = [
    { field: "service", value: "unknown" },
    { field: "address", value: " " },
    { field: "address", value: "я".repeat(101) },
    { field: "date", value: "2026-10-20" },
    { field: "date", value: "2027-02-29" },
    { field: "date", value: "2027-13-01" },
  ];

  for (const { field, value } of invalidValues) {
    const errors = validateRequest(
      { ...input, [field]: value },
      services,
      today,
    );
    assert.ok(errors[field], `${field}: ${value}`);
  }
});
