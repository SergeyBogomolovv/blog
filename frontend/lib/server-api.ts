// Доступ к Go API только с сервера Next.js; токен не попадает в клиентские props.
import "server-only";

import { cookies } from "next/headers";
import { notFound } from "next/navigation";
import { cache } from "react";

import type { User } from "@/lib/api";

const backendURL = process.env.BACKEND_URL || "http://127.0.0.1:8080";
const appOrigin = process.env.APP_ORIGIN || "http://localhost:3000";

// ApiError отделяет ожидаемый отказ API от ошибки рендера приложения.
export class ApiError extends Error {
  // Сохраняет сообщение для пользователя и HTTP-статус ответа API.
  constructor(
    message: string,
    public status: number,
  ) {
    super(message);
  }
}

// backendRequest пересылает только cookie авторизации, ограничивая время ожидания API.
async function backendRequest(
  path: string,
  method: string,
  data?: unknown,
  authenticated = true,
) {
  const token = authenticated
    ? (await cookies()).get("access_token")?.value
    : undefined;

  try {
    return await fetch(`${backendURL}/api${path}`, {
      method,
      cache: "no-store",
      signal: AbortSignal.timeout(10000),
      headers: {
        "Content-Type": "application/json",
        "X-Requested-With": "studyblog",
        Origin: appOrigin,
        ...(token
          ? { Cookie: `access_token=${encodeURIComponent(token)}` }
          : {}),
      },
      body: data === undefined ? undefined : JSON.stringify(data),
    });
  } catch {
    throw new ApiError(
      "Не удалось связаться с сервером. Проверьте, что бэкенд запущен.",
      503,
    );
  }
}

// result извлекает публичный ответ, не раскрывая внутренние ошибки сервера.
async function result<T>(response: Response): Promise<T> {
  const body = await response.json().catch(() => ({}));

  if (!response.ok) {
    throw new ApiError(
      body.error || "Не удалось выполнить запрос. Попробуйте позже.",
      response.status,
    );
  }

  return body as T;
}

// readAPI читает данные без общего кеша между пользователями.
export async function readAPI<T>(
  path: string,
  authenticated = true,
): Promise<T> {
  return result<T>(await backendRequest(path, "GET", undefined, authenticated));
}

// readPublic читает общедоступные данные; отсутствующий ресурс отдаёт страницу 404.
export async function readPublic<T>(path: string): Promise<T> {
  try {
    return await readAPI<T>(path, false);
  } catch (error) {
    if (error instanceof ApiError && error.status === 404) {
      notFound();
    }

    throw error;
  }
}

// getSession проверяет аккаунт один раз в рамках серверного рендера, а не между запросами.
export const getSession = cache(
  async function getSession(): Promise<User | null> {
    if (!(await cookies()).has("access_token")) {
      return null;
    }

    try {
      return (await readAPI<{ user: User | null }>("/auth/me")).user;
    } catch (error) {
      if (error instanceof ApiError && error.status === 401) {
        return null;
      }

      throw error;
    }
  },
);

// mutateAPI вызывается только из Server Actions и передаёт Set-Cookie в ответ браузеру.
export async function mutateAPI<T>(
  path: string,
  method: string,
  data?: unknown,
): Promise<T> {
  const response = await backendRequest(
    path,
    method,
    data,
    path !== "/auth/login" && path !== "/auth/register",
  );
  const header = response.headers
    .getSetCookie()
    .find((value) => value.startsWith("access_token="));

  if (header) {
    const attributes = header.split(";").map((value) => value.trim());
    const value = attributes[0].slice("access_token=".length);
    const maxAge = Number(
      attributes
        .find((value) => value.toLowerCase().startsWith("max-age="))
        ?.split("=")[1] || "7200",
    );
    (await cookies()).set("access_token", value, {
      httpOnly: true,
      sameSite: "strict",
      path: "/",
      secure: attributes.some((value) => value.toLowerCase() === "secure"),
      maxAge,
    });
  }

  return result<T>(response);
}
