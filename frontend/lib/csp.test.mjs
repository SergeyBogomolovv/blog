// Проверка разрешений ресурсов и изоляции настроек DevTools от собранного сайта.
import assert from "node:assert/strict";
import { test } from "node:test";

import { contentSecurityPolicy } from "./csp.ts";

/** Разбирает директивы CSP для проверки отдельных типов ресурсов. */
function directives(policy) {
  return new Map(
    policy.split("; ").map((directive) => {
      const [name, ...sources] = directive.split(" ");

      return [name, sources];
    }),
  );
}

test("собранный сайт требует nonce и запрещает небезопасные разрешения", () => {
  const policy = contentSecurityPolicy("production-nonce");
  const rules = directives(policy);

  assert.deepEqual(rules.get("script-src"), [
    "'self'",
    "'nonce-production-nonce'",
  ]);
  assert.deepEqual(rules.get("style-src"), [
    "'self'",
    "'nonce-production-nonce'",
  ]);
  assert.equal(rules.has("style-src-elem"), false);
  assert.equal(rules.has("style-src-attr"), false);
  assert.doesNotMatch(policy, /'unsafe-inline'|'unsafe-eval'|\*/);
  assert.equal(policy, contentSecurityPolicy("production-nonce", false));
});

test("разработка разрешает теги стилей DevTools, сохраняя защиту других ресурсов", () => {
  const rules = directives(contentSecurityPolicy("development-nonce", true));

  assert.deepEqual(rules.get("style-src-elem"), ["'self'", "'unsafe-inline'"]);
  assert.deepEqual(rules.get("style-src"), [
    "'self'",
    "'nonce-development-nonce'",
  ]);
  assert.equal(rules.has("style-src-attr"), false);
  assert.deepEqual(rules.get("script-src"), [
    "'self'",
    "'nonce-development-nonce'",
    "'unsafe-eval'",
  ]);

  for (const development of [true, false]) {
    const policy = directives(contentSecurityPolicy("nonce", development));
    assert.deepEqual(policy.get("object-src"), ["'none'"]);
    assert.deepEqual(policy.get("frame-ancestors"), ["'none'"]);
    assert.deepEqual(policy.get("form-action"), ["'self'"]);
    assert.deepEqual(policy.get("connect-src"), ["'self'"]);
  }
});
