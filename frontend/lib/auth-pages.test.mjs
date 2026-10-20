// Регрессия переходов между формами авторизации и личным кабинетом.
import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import { test } from "node:test";
import { runInNewContext } from "node:vm";
import ts from "typescript";

/** Загружает реальную страницу, подменяя только сессию, навигацию и дочерние компоненты. */
function loadPage(route, session) {
  const file = new URL(`../app/${route}/page.tsx`, import.meta.url);
  const compiled = ts.transpileModule(readFileSync(file, "utf8"), {
    compilerOptions: {
      module: ts.ModuleKind.CommonJS,
      jsx: ts.JsxEmit.ReactJSX,
      target: ts.ScriptTarget.ES2022,
    },
  });
  const pageModule = { exports: {} };
  const dependencies = {
    "react/jsx-runtime": {
      jsx: (type, props) => ({ type, props }),
      jsxs: (type, props) => ({ type, props }),
    },
    "next/navigation": {
      redirect: (destination) => {
        throw Object.assign(new Error("redirect"), { destination });
      },
    },
    "@/lib/server-api": {
      getSession: async () => session,
      readAPI: async () => ({ services: [], requests: [] }),
    },
    "@/components/auth-form": { AuthForm: "AuthForm" },
    "@/components/auth-layout": { AuthLayout: "AuthLayout" },
    "@/components/dashboard": { Dashboard: "Dashboard" },
    "@/components/service-unavailable": {
      ServiceUnavailable: "ServiceUnavailable",
    },
  };

  runInNewContext(compiled.outputText, {
    module: pageModule,
    exports: pageModule.exports,
    require: (name) => {
      assert.ok(
        Object.hasOwn(dependencies, name),
        `Неизвестный импорт: ${name}`,
      );

      return dependencies[name];
    },
  });

  return pageModule.exports.default;
}

for (const route of ["login", "register"]) {
  test(`${route}: гость остаётся на форме после отказа в проверке сессии`, async () => {
    const page = loadPage(route, { user: null, unavailable: false });
    const rendered = await page({
      searchParams: Promise.resolve({ service: "deep", registered: "1" }),
    });

    assert.equal(rendered.type, "AuthLayout");
    assert.equal(rendered.props.children.type, "AuthForm");
    assert.equal(rendered.props.children.props.mode, route);
    assert.equal(rendered.props.children.props.service, "deep");
  });

  test(`${route}: только действующий пользователь перенаправляется в кабинет`, async () => {
    const page = loadPage(route, { user: { id: 1 }, unavailable: false });

    await assert.rejects(
      page({ searchParams: Promise.resolve({ service: "deep" }) }),
      (error) => error.destination === "/dashboard?service=deep",
    );
  });

  test(`${route}: сбой API показывает сообщение без перенаправления`, async () => {
    const page = loadPage(route, { user: null, unavailable: true });
    const rendered = await page({ searchParams: Promise.resolve({}) });

    assert.equal(rendered.type, "ServiceUnavailable");
  });
}

test("кабинет гостя переходит на вход, где цепочка перенаправлений заканчивается", async () => {
  const session = { user: null, unavailable: false };
  const dashboard = loadPage("dashboard", session);
  const login = loadPage("login", session);

  await assert.rejects(
    dashboard({ searchParams: Promise.resolve({ service: "deep" }) }),
    (error) => error.destination === "/login?service=deep",
  );

  const rendered = await login({
    searchParams: Promise.resolve({ service: "deep" }),
  });
  assert.equal(rendered.type, "AuthLayout");
});
