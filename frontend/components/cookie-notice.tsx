// Уведомление о необходимых cookie и сохранение подтверждения.
"use client";

import Link from "next/link";
import { useState } from "react";
import { Cookie } from "lucide-react";

import { api, errorMessage } from "@/lib/api";

/** Показывает уведомление до подтверждения пользователем. */
export function CookieNotice() {
  const [visible, setVisible] = useState(true);
  const [pending, setPending] = useState(false);
  const [error, setError] = useState("");

  /** Сохраняет подтверждение на сервере и скрывает уведомление. */
  async function acceptNotice() {
    setPending(true);
    setError("");

    try {
      await api("/cookie-notice", {});
      setVisible(false);
    } catch (error) {
      setError(errorMessage(error));
    } finally {
      setPending(false);
    }
  }

  if (!visible) {
    return null;
  }

  return (
    <aside
      className="cookie-notice"
      aria-label="Использование cookie"
    >
      <Cookie size={25} />
      <div>
        <strong>Немного cookie для порядка</strong>
        <p>
          Сохраняем вход и защищаем ваши формы.{" "}
          <Link href="/privacy">Подробнее</Link>
        </p>
        {error && (
          <p
            className="field-error"
            role="alert"
          >
            {error}
          </p>
        )}
      </div>
      <button
        type="button"
        className="button button-small"
        disabled={pending}
        onClick={acceptNotice}
      >
        {pending ? "Сохраняем…" : "Понятно"}
      </button>
    </aside>
  );
}
