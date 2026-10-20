"use client";

// Формы авторизации отправляются через Server Actions; JWT не доступен компоненту.
import Link from "next/link";
import { useActionState } from "react";

import { ArrowLeft, ArrowRight } from "lucide-react";

import { loginAction, registerAction } from "@/app/actions";
import { Button } from "@/components/ui/button";
import { CSRFField } from "@/components/csrf";

// AuthForm отображает состояние действия и сохраняет несекретные поля при ошибке.
export function AuthForm({ register = false }: { register?: boolean }) {
  const [state, formAction, busy] = useActionState(
    register ? registerAction : loginAction,
    {},
    register ? "/register" : "/login",
  );

  return (
    <div className="form-layout">
      <Link
        href="/"
        className="back-link"
      >
        <ArrowLeft size={16} />К историям
      </Link>
      <div className="form-card">
        <span className="eyebrow">Между строк</span>
        <h1>
          {register ? "Ваша история начинается здесь" : "Рады видеть вас снова"}
        </h1>
        <p className="form-intro">
          {register
            ? "Создайте аккаунт, чтобы писать и участвовать в обсуждениях."
            : "Войдите, чтобы продолжить писать свою историю."}
        </p>
        <form
          action={formAction}
          autoComplete="off"
        >
          <CSRFField />
          {register && (
            <>
              <label htmlFor="name">Как к вам обращаться</label>
              <input
                id="name"
                name="name"
                defaultValue={state.values?.name ?? ""}
                required
                minLength={2}
                maxLength={60}
                placeholder="Ваше имя"
                autoComplete="off"
              />
              <label htmlFor="email">Электронная почта</label>
              <input
                id="email"
                name="email"
                defaultValue={state.values?.email ?? ""}
                type="email"
                required
                maxLength={254}
                placeholder="you@example.com"
                autoComplete="off"
              />
              <span className="field-hint">Почта видна только вам.</span>
            </>
          )}
          <label htmlFor="login">Логин</label>
          <input
            id="login"
            name="login"
            defaultValue={state.values?.login ?? ""}
            required
            pattern="[a-z0-9_]{3,24}"
            minLength={3}
            maxLength={24}
            autoComplete="off"
            aria-describedby="login-hint"
            placeholder="your_name"
          />
          <span
            id="login-hint"
            className="field-hint"
          >
            3–24 символа: маленькие латинские буквы, цифры и _.
          </span>
          <label htmlFor="password">Пароль</label>
          <input
            id="password"
            name="password"
            type="password"
            required
            minLength={10}
            maxLength={72}
            autoComplete="off"
            aria-describedby="password-hint"
            placeholder="Не меньше 10 символов"
          />
          <span
            id="password-hint"
            className="field-hint"
          >
            Не менее 10 символов.
          </span>
          {register && (
            <label className="checkbox-label">
              <input
                name="consent"
                type="checkbox"
                required
              />
              <span>
                Согласие на обработку данных согласно{" "}
                <Link href="/privacy">политике обработки данных</Link>
              </span>
            </label>
          )}
          {state.error && (
            <p
              className="error-box"
              role="alert"
            >
              {state.error}
            </p>
          )}
          <Button
            type="submit"
            className="wide-button"
            disabled={busy}
          >
            {busy ? "Подождите…" : register ? "Создать аккаунт" : "Войти"}
            <ArrowRight size={17} />
          </Button>
        </form>
        <p className="form-switch">
          {register ? "Уже есть аккаунт?" : "Ещё нет аккаунта?"}{" "}
          <Link href={register ? "/login" : "/register"}>
            {register ? "Войти" : "Присоединиться"}
          </Link>
        </p>
      </div>
    </div>
  );
}
