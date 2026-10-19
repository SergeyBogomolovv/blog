"use server";

// Изменения данных выполняются на сервере; Go повторно проверяет JWT, поля и права.
import { revalidatePath } from "next/cache";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";

import type { FormState } from "@/lib/form-state";
import { ApiError, getSession, mutateAPI } from "@/lib/server-api";

// text получает только строковое значение и не принимает файлы вместо полей.
function text(data: FormData, name: string): string {
  const value = data.get(name);

  return typeof value === "string" ? value : "";
}

// failure сохраняет несекретные поля для повторного ввода и показывает безопасную ошибку.
function failure(error: unknown, values?: Record<string, string>): FormState {
  return {
    error:
      error instanceof ApiError
        ? error.message
        : "Не удалось выполнить действие. Попробуйте позже.",
    values,
  };
}

// requireUser повторно проверяет сессию при вызове действия напрямую.
async function requireUser(admin = false) {
  const user = await getSession();

  if (!user) {
    throw new ApiError("Войдите в аккаунт заново", 401);
  }

  if (admin && user.role !== "admin") {
    throw new ApiError("Недостаточно прав для этого действия", 403);
  }

  return user;
}

// identifier проверяет идентификатор из недоверенного поля формы.
function identifier(data: FormData, name: string): string {
  const value = text(data, name);

  if (!/^[1-9]\d{0,15}$/.test(value)) {
    throw new ApiError("Некорректный идентификатор", 400);
  }

  return value;
}

// loginAction входит в аккаунт и обновляет серверную навигацию после установки cookie.
export async function loginAction(
  _state: FormState,
  data: FormData,
): Promise<FormState> {
  const login = text(data, "login");

  try {
    await mutateAPI("/auth/login", "POST", {
      login,
      password: text(data, "password"),
    });
  } catch (error) {
    return failure(error, { login });
  }

  revalidatePath("/", "layout");
  redirect("/");
}

// registerAction проверяет согласие до передачи регистрационных данных в Go API.
export async function registerAction(
  _state: FormState,
  data: FormData,
): Promise<FormState> {
  const values = {
    login: text(data, "login"),
    name: text(data, "name"),
    email: text(data, "email"),
  };

  if (data.get("consent") !== "on") {
    return { error: "Подтвердите согласие на обработку данных.", values };
  }

  try {
    await mutateAPI("/auth/register", "POST", {
      ...values,
      password: text(data, "password"),
      consent: true,
    });
  } catch (error) {
    return failure(error, values);
  }

  revalidatePath("/", "layout");
  redirect("/");
}

// logoutAction отзывает токены в Go и удаляет cookie даже при уже истёкшей сессии.
export async function logoutAction(): Promise<FormState> {
  try {
    await mutateAPI("/auth/logout", "POST", {});
  } catch (error) {
    if (!(error instanceof ApiError && error.status === 401)) {
      return failure(error);
    }
  }

  (await cookies()).delete("access_token");
  revalidatePath("/", "layout");
  redirect("/");
}

// createPostAction публикует текст от имени проверенного пользователя.
export async function createPostAction(
  _state: FormState,
  data: FormData,
): Promise<FormState> {
  const values = { title: text(data, "title"), body: text(data, "body") };
  let id: number;

  try {
    await requireUser();
    id = (await mutateAPI<{ id: number }>("/posts", "POST", values)).id;
  } catch (error) {
    return failure(error, values);
  }

  revalidatePath("/", "layout");
  redirect(`/posts/${id}`);
}

// commentAction добавляет комментарий и возвращает обновлённое обсуждение в том же ответе.
export async function commentAction(
  _state: FormState,
  data: FormData,
): Promise<FormState> {
  const body = text(data, "body");
  let id: string;

  try {
    await requireUser();
    id = identifier(data, "post_id");
    await mutateAPI(`/posts/${id}/comments`, "POST", { body });
  } catch (error) {
    return failure(error, { body });
  }

  revalidatePath(`/posts/${id}`);
  revalidatePath("/");

  return { success: "Комментарий добавлен.", values: { body: "" } };
}

// profileAction меняет только профиль владельца сессии и обновляет его публикации.
export async function profileAction(
  _state: FormState,
  data: FormData,
): Promise<FormState> {
  const values = { name: text(data, "name"), bio: text(data, "bio") };

  try {
    await requireUser();
    await mutateAPI("/profile", "PATCH", values);
  } catch (error) {
    return failure(error, values);
  }

  revalidatePath("/", "layout");

  return { success: "Профиль обновлён.", values };
}

// deletePostAction требует явного подтверждения; проверка автора остаётся в Go-сервисе.
export async function deletePostAction(
  _state: FormState,
  data: FormData,
): Promise<FormState> {
  try {
    await requireUser();

    if (data.get("confirm") !== "on") {
      throw new ApiError("Подтвердите удаление", 400);
    }

    const id = identifier(data, "post_id");
    await mutateAPI(`/posts/${id}`, "DELETE");
  } catch (error) {
    return failure(error);
  }

  revalidatePath("/", "layout");
  redirect("/");
}

// banAction проверяет роль администратора и передаёт явное состояние блокировки.
export async function banAction(
  _state: FormState,
  data: FormData,
): Promise<FormState> {
  try {
    await requireUser(true);
    const id = identifier(data, "user_id");
    const banned = text(data, "banned");

    if (!["true", "false"].includes(banned) || data.get("confirm") !== "on") {
      throw new ApiError("Подтвердите действие", 400);
    }

    await mutateAPI(`/admin/users/${id}/ban`, "PATCH", {
      banned: banned === "true",
    });
  } catch (error) {
    return failure(error);
  }

  revalidatePath("/", "layout");

  return { success: "Статус пользователя обновлён." };
}

// acceptCookiesAction запоминает прочтение уведомления в отдельной необходимой cookie.
export async function acceptCookiesAction(
  _state: FormState,
  data: FormData,
): Promise<FormState> {
  (await cookies()).set("cookie_notice", "accepted", {
    httpOnly: true,
    sameSite: "strict",
    path: "/",
    maxAge: 31536000,
    secure: process.env.COOKIE_SECURE === "true",
  });
  revalidatePath("/", "layout");

  if (text(data, "destination") === "privacy") {
    redirect("/privacy");
  }

  return {};
}
