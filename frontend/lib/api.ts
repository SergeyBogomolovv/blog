// Типы данных и обращения браузера к Go API через Next.js.
export interface User {
  id: number;
  surname: string;
  name: string;
  patronymic: string;
  login: string;
  phone: string;
  email: string;
  consent_at: string;
  consent_version: string;
}

/** Структура данных CleaningService. */
export interface CleaningService {
  id: string;
  name: string;
  description: string;
}

/** Структура данных CleaningRequest. */
export interface CleaningRequest {
  id: number;
  service: string;
  address: string;
  cleaning_date: string;
  created_at: string;
}

/** Структура данных RegistrationValues. */
export interface RegistrationValues {
  surname: string;
  name: string;
  patronymic: string;
  login: string;
  password: string;
  phone: string;
  email: string;
  consent: boolean;
}

/** Структура данных RequestValues. */
export interface RequestValues {
  service: string;
  address: string;
  date: string;
}

/** Структура данных FieldErrors. */
export type FieldErrors = Record<string, string>;

/** Структура данных ErrorBody. */
interface ErrorBody {
  error?: string;
  fields?: FieldErrors;
}

/** Ошибка API с безопасным сообщением и ошибками полей. */
export class ApiError extends Error {
  readonly status: number;
  readonly fields: FieldErrors;

  /** Создаёт ошибку API с HTTP-статусом и ошибками полей. */
  constructor(message: string, status: number, fields: FieldErrors = {}) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.fields = fields;
  }
}

/** Возвращает понятное сообщение для ошибки запроса. */
export function errorMessage(error: unknown): string {
  if (error instanceof Error) {
    return error.message;
  }

  return "Не удалось выполнить запрос. Попробуйте ещё раз.";
}

/** Форматирует дату для отображения по московскому времени. */
export function date(value: string): string {
  return new Intl.DateTimeFormat("ru-RU", {
    day: "numeric",
    month: "long",
    year: "numeric",
    timeZone: "Europe/Moscow",
  }).format(new Date(value));
}

// Календарный день определяется по московскому времени, а не по UTC.
export function moscowToday(): string {
  return new Intl.DateTimeFormat("en-CA", {
    timeZone: "Europe/Moscow",
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  }).format(new Date());
}

/** Получает подписанный CSRF-токен перед изменением данных. */
async function csrfToken(): Promise<string> {
  const response = await fetch("/api/csrf", {
    cache: "no-store",
    signal: AbortSignal.timeout(12000),
  });

  if (!response.ok) {
    throw new ApiError(
      "Не удалось подготовить форму. Повторите попытку.",
      response.status,
    );
  }

  const body: { token: string } = await response.json();

  return body.token;
}

/** Отправляет запрос через общий адрес сайта и обрабатывает ошибки. */
export async function api<T>(path: string, data?: unknown): Promise<T> {
  const isMutation = data !== undefined;

  try {
    const headers = new Headers({
      "Content-Type": "application/json",
      "X-Requested-With": "cleaning",
    });

    if (isMutation) {
      headers.set("X-CSRF-Token", await csrfToken());
    }

    const response = await fetch(`/api${path}`, {
      method: isMutation ? "POST" : "GET",
      cache: "no-store",
      headers,
      body: isMutation ? JSON.stringify(data) : undefined,
      signal: AbortSignal.timeout(12000),
    });

    if (!response.ok) {
      const body: ErrorBody = await response.json();

      throw new ApiError(
        body.error || "Не удалось выполнить запрос.",
        response.status,
        body.fields || {},
      );
    }

    const body: T = await response.json();

    return body;
  } catch (error) {
    if (error instanceof ApiError) {
      throw error;
    }

    throw new ApiError(
      "Не удалось связаться с сервером. Попробуйте ещё раз.",
      503,
    );
  }
}
