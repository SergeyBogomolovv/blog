// Тестовый ключ не связан с .env; тесты не требуют запущенного API.
import assert from "node:assert/strict";
import { test } from "node:test";
import { createCSRFToken, sameCSRFToken, validCSRFToken } from "./csrf.ts";

process.env.JWT_SECRET = "01".repeat(32);
const vector =
  "0".repeat(64) +
  ".2f897b9659eb48931898af16a1c0abbc3e58d2f5710d194a86463ee26e7f463c";

test("подпись совместима с эталоном Go", () => {
  assert.equal(validCSRFToken(vector, "session-a"), true);
  assert.equal(validCSRFToken(vector, "session-b"), false);
  assert.equal(validCSRFToken(vector), false);
});

test("случайные токены действуют только в своей сессии", () => {
  const token = createCSRFToken("session-a");
  assert.match(token, /^[a-f0-9]{64}\.[a-f0-9]{64}$/);
  assert.notEqual(token, createCSRFToken("session-a"));
  assert.equal(validCSRFToken(token, "session-a"), true);
  assert.equal(validCSRFToken(token, "session-b"), false);
  assert.equal(validCSRFToken(createCSRFToken()), true);
});

test("отсутствующие и поддельные токены отклоняются", () => {
  for (const token of [
    "",
    "a".repeat(129),
    "1" + vector.slice(1),
    vector.slice(0, -1) + "0",
    vector + ".extra",
  ]) {
    assert.equal(validCSRFToken(token, "session-a"), false);
  }

  assert.equal(sameCSRFToken(vector, vector), true);
  assert.equal(sameCSRFToken(vector, "1" + vector.slice(1)), false);
  assert.equal(sameCSRFToken("", ""), false);
  assert.equal(sameCSRFToken(vector, "я".repeat(129)), false);
});

test("смена серверного ключа отзывает токен, отсутствие ключа запрещает выдачу", () => {
  const original = process.env.JWT_SECRET;

  try {
    process.env.JWT_SECRET = "02".repeat(32);
    assert.equal(validCSRFToken(vector, "session-a"), false);
    delete process.env.JWT_SECRET;
    assert.throws(() => createCSRFToken());
    assert.throws(() => validCSRFToken(vector, "session-a"));
  } finally {
    process.env.JWT_SECRET = original;
  }
});
