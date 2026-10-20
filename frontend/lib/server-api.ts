// Обращения серверных компонентов к API с передачей cookie.
import "server-only";

import { cookies } from "next/headers";
import { cache } from "react";

import { ApiError, type User } from "@/lib/api";

export const backendURL = process.env.BACKEND_URL || "http://127.0.0.1:8080";

/** Получает данные API без кеширования и с cookie текущего запроса. */
export async function readAPI<T>(path: string): Promise<T> {
  const cookieJar = await cookies();
  let response: Response;

  try {
    response = await fetch(`${backendURL}/api${path}`, {
      cache: "no-store",
      headers: { Cookie: cookieJar.toString() },
      signal: AbortSignal.timeout(5000),
    });
  } catch {
    throw new ApiError(
      "Сервис временно недоступен. Проверьте запуск Go API и базы данных.",
      503,
    );
  }

  if (!response.ok) {
    const body: { error?: string } = await response.json();

    throw new ApiError(
      body.error || "Не удалось загрузить данные.",
      response.status,
    );
  }

  const body: T = await response.json();

  return body;
}

/** Состояние сессии отличает гостя от временной недоступности сервиса. */
interface SessionState {
  user: User | null;
  unavailable: boolean;
}

/** Читает сессию один раз за рендер, сохраняя сбои API как явное состояние. */
export const getSession = cache(
  async function getSession(): Promise<SessionState> {
    const cookieJar = await cookies();

    if (!cookieJar.has("cleaning.sid")) {
      return { user: null, unavailable: false };
    }

    try {
      const response = await readAPI<{ user: User }>("/auth/me");

      return { user: response.user, unavailable: false };
    } catch (error) {
      if (error instanceof ApiError && error.status === 401) {
        return { user: null, unavailable: false };
      }

      return { user: null, unavailable: true };
    }
  },
);
