// Регистрация и вход с проверкой полей до отправки.
"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { type FormEvent, useState } from "react";
import { ArrowUpRight, LockKeyhole } from "lucide-react";

import {
  AuthTextField,
  loginField,
  PasswordField,
  registrationFields,
} from "@/components/auth-fields";
import {
  api,
  ApiError,
  errorMessage,
  /** Структура данных FieldErrors. */
  type FieldErrors,
  /** Структура данных RegistrationValues. */
  type RegistrationValues,
} from "@/lib/api";
import { validateLogin, validateRegistration } from "@/lib/validation";

/** Параметры компонента AuthFormProps. */
interface AuthFormProps {
  mode: "login" | "register";
  registered?: boolean;
  service?: string;
}

/** Читает поля формы, сохраняя пароль без изменения пробелов. */
function registrationInput(data: FormData): RegistrationValues {
  /** Получает строковое поле и удаляет пробелы по краям. */
  function text(name: string): string {
    return String(data.get(name) || "").trim();
  }

  return {
    surname: text("surname"),
    name: text("name"),
    patronymic: text("patronymic"),
    login: text("login"),
    phone: text("phone"),
    email: text("email"),
    password: String(data.get("password") || ""),
    consent: data.has("consent"),
  };
}

/** Формирует адрес перехода с выбранной услугой. */
function destinationAfterSubmit(
  mode: AuthFormProps["mode"],
  service: string,
): string {
  const parameters = new URLSearchParams();
  const pathname = mode === "register" ? "/login" : "/dashboard";

  if (mode === "register") {
    parameters.set("registered", "1");
  }

  if (service) {
    parameters.set("service", service);
  }

  const query = parameters.toString();

  return query ? `${pathname}?${query}` : pathname;
}

/** Проверяет и отправляет данные регистрации либо входа. */
export function AuthForm({
  mode,
  registered = false,
  service = "",
}: AuthFormProps) {
  const router = useRouter();
  const [errors, setErrors] = useState<FieldErrors>({});
  const [error, setError] = useState("");
  const [pending, setPending] = useState(false);

  const isRegistration = mode === "register";
  const fields = isRegistration ? registrationFields : [loginField];
  const submitLabel = isRegistration ? "Создать аккаунт" : "Войти в кабинет";
  const switchPath = isRegistration ? "/login" : "/register";
  const serviceQuery = service ? `?service=${encodeURIComponent(service)}` : "";

  /** Проверяет поля и отправляет данные авторизации. */
  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");

    const form = event.currentTarget;
    const input = registrationInput(new FormData(form));
    const validationErrors = isRegistration
      ? validateRegistration(input)
      : validateLogin(input);
    const invalidFields = Object.keys(validationErrors);
    setErrors(validationErrors);

    if (invalidFields.length > 0) {
      const firstField = form.querySelector<HTMLInputElement>(
        `[name="${invalidFields[0]}"]`,
      );
      firstField?.focus();

      return;
    }

    setPending(true);

    try {
      const values = isRegistration
        ? input
        : { login: input.login, password: input.password };

      await api(`/auth/${mode}`, values);
      router.push(destinationAfterSubmit(mode, service));
      router.refresh();
    } catch (error) {
      if (error instanceof ApiError) {
        setErrors(error.fields);
      }

      setError(errorMessage(error));
      setPending(false);
    }
  }

  return (
    <form
      method="post"
      className="auth-form"
      onSubmit={submit}
      autoComplete="off"
    >
      {registered && (
        <div
          className="success-box"
          role="status"
        >
          Аккаунт создан. Войдите со своим логином и паролем.
        </div>
      )}
      {error && (
        <div
          className="error-box"
          role="alert"
        >
          {error}
        </div>
      )}
      <div className={isRegistration ? "form-grid" : "form-stack"}>
        {fields.map((field) => (
          <AuthTextField
            key={field.name}
            field={field}
            disabled={pending}
            error={errors[field.name]}
          />
        ))}
        <PasswordField
          registration={isRegistration}
          disabled={pending}
          error={errors.password}
        />
      </div>
      {isRegistration && (
        <div className="consent-field">
          <label>
            <input
              type="checkbox"
              name="consent"
              required
              disabled={pending}
              aria-describedby={errors.consent ? "consent-error" : undefined}
            />
            <span>
              Я даю{" "}
              <Link
                href="/privacy"
                target="_blank"
                rel="noopener"
              >
                согласие на обработку персональных данных
              </Link>
            </span>
          </label>
          {errors.consent && (
            <p
              className="field-error"
              id="consent-error"
            >
              {errors.consent}
            </p>
          )}
        </div>
      )}
      <button
        type="submit"
        className="button wide-button"
        disabled={pending}
      >
        {pending ? "Подождите…" : submitLabel}
        <ArrowUpRight size={18} />
      </button>
      <p className="auth-switch">
        {isRegistration ? "Уже есть аккаунт?" : "Пока нет аккаунта?"}{" "}
        <Link href={`${switchPath}${serviceQuery}`}>
          {isRegistration ? "Войти" : "Зарегистрироваться"}
        </Link>
      </p>
      <p className="form-security">
        <LockKeyhole size={13} />
        {isRegistration
          ? "Ваш пароль будет храниться только в виде хеша"
          : "Ваши заявки доступны только вам"}
      </p>
    </form>
  );
}
