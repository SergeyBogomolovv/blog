"use client";

import { type FocusEvent, useState } from "react";
import { Eye, EyeOff } from "lucide-react";

import { passwordHint } from "@/lib/validation";

type TextFieldName = "surname" | "name" | "patronymic" | "login" | "phone" | "email";

interface TextFieldDefinition {
  name: TextFieldName;
  label: string;
  placeholder: string;
  hint: string;
  type: "text" | "tel" | "email";
}

export const loginField: TextFieldDefinition = {
  name: "login",
  label: "Логин",
  placeholder: "Ivan",
  hint: "Только латинские буквы, до 20 символов",
  type: "text",
};

export const registrationFields: TextFieldDefinition[] = [
  {
    name: "surname",
    label: "Фамилия",
    placeholder: "Иванов",
    hint: "Русские буквы, первая заглавная",
    type: "text",
  },
  {
    name: "name",
    label: "Имя",
    placeholder: "Иван",
    hint: "Можно двойное: Анна-Мария",
    type: "text",
  },
  {
    name: "patronymic",
    label: "Отчество",
    placeholder: "Иванович",
    hint: "Допустимо: Рашид оглы",
    type: "text",
  },
  loginField,
  {
    name: "phone",
    label: "Телефон",
    placeholder: "8(999)123-45-67",
    hint: "Строго в формате 8(999)123-45-67",
    type: "tel",
  },
  {
    name: "email",
    label: "Электронная почта",
    placeholder: "name@mail.ru",
    hint: "До 20 символов, включая домен",
    type: "email",
  },
];

interface AuthTextFieldProps {
  field: TextFieldDefinition;
  disabled: boolean;
  error?: string;
}

function trimInput(event: FocusEvent<HTMLInputElement>) {
  event.currentTarget.value = event.currentTarget.value.trim();
}

export function AuthTextField({ field, disabled, error }: AuthTextFieldProps) {
  const hintID = `${field.name}-hint`;
  const errorID = `${field.name}-error`;
  const description = error ? `${hintID} ${errorID}` : hintID;

  return (
    <div className="field">
      <label htmlFor={field.name}>{field.label}</label>
      <input
        id={field.name}
        name={field.name}
        type={field.type}
        placeholder={field.placeholder}
        required
        maxLength={20}
        disabled={disabled}
        autoCapitalize="none"
        spellCheck={false}
        onBlur={trimInput}
        aria-invalid={Boolean(error)}
        aria-describedby={description}
      />
      <span className="field-hint" id={hintID}>{field.hint}</span>
      {error && <span className="field-error" id={errorID}>{error}</span>}
    </div>
  );
}

interface PasswordFieldProps {
  registration: boolean;
  disabled: boolean;
  error?: string;
}

export function PasswordField({ registration, disabled, error }: PasswordFieldProps) {
  const [visible, setVisible] = useState(false);
  const placeholder = registration
    ? "Придумайте надёжный пароль"
    : "Введите пароль";
  const hint = registration
    ? passwordHint
    : "Введите пароль, указанный при регистрации.";

  return (
    <div className={registration ? "field full-width" : "field"}>
      <label htmlFor="password">Пароль</label>
      <div className="password-input">
        <input
          id="password"
          name="password"
          type={visible ? "text" : "password"}
          minLength={registration ? 8 : undefined}
          placeholder={placeholder}
          required
          disabled={disabled}
          aria-invalid={Boolean(error)}
          aria-describedby="password-hint password-error"
        />
        <button
          type="button"
          className="icon-button"
          aria-label={visible ? "Скрыть пароль" : "Показать пароль"}
          aria-pressed={visible}
          onClick={() => setVisible((current) => !current)}
        >
          {visible ? <EyeOff size={18} /> : <Eye size={18} />}
        </button>
      </div>
      <span className="field-hint" id="password-hint">{hint}</span>
      <span className="field-error" id="password-error">{error}</span>
    </div>
  );
}
