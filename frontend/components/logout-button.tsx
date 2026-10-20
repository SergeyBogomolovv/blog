// Выход из аккаунта с отзывом сессии на сервере.
"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";

import { api, errorMessage } from "@/lib/api";

/** Отображает выход и состояние выполнения запроса. */
export function LogoutButton() {
  const router = useRouter();
  const [pending, setPending] = useState(false);
  const [error, setError] = useState("");

  /** Отзывает сессию и возвращает пользователя к форме входа. */
  async function logout() {
    setPending(true);
    setError("");

    try {
      await api("/auth/logout", {});
      router.push("/");
      router.refresh();
    } catch (error) {
      setError(errorMessage(error));
      setPending(false);
    }
  }

  return (
    <div className="logout">
      <button
        type="button"
        className="login-link"
        disabled={pending}
        onClick={logout}
      >
        {pending ? "Выходим…" : "Выйти"}
      </button>
      {error && (
        <p
          className="field-error"
          role="alert"
        >
          {error}
        </p>
      )}
    </div>
  );
}
