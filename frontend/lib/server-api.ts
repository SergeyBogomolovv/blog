import "server-only";

import { cookies } from "next/headers";
import { cache } from "react";

import { ApiError, type User } from "@/lib/api";

export const backendURL = process.env.BACKEND_URL || "http://127.0.0.1:8080";

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

export const getSession = cache(async function getSession(): Promise<User | null> {
  const cookieJar = await cookies();

  if (!cookieJar.has("cleaning.sid")) {
    return null;
  }

  try {
    const response = await readAPI<{ user: User }>("/auth/me");

    return response.user;
  } catch (error) {
    if (error instanceof ApiError && error.status === 401) {
      return null;
    }

    throw error;
  }
});
