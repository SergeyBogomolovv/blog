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

export interface CleaningService {
  id: string;
  name: string;
  description: string;
}

export interface CleaningRequest {
  id: number;
  service: string;
  address: string;
  cleaning_date: string;
  created_at: string;
}

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

export interface RequestValues {
  service: string;
  address: string;
  date: string;
}

export type FieldErrors = Record<string, string>;

interface ErrorBody {
  error?: string;
  fields?: FieldErrors;
}

export class ApiError extends Error {
  readonly status: number;
  readonly fields: FieldErrors;

  constructor(message: string, status: number, fields: FieldErrors = {}) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.fields = fields;
  }
}

export function errorMessage(error: unknown): string {
  if (error instanceof Error) {
    return error.message;
  }

  return "Не удалось выполнить запрос. Попробуйте ещё раз.";
}

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
