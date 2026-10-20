// Браузер обращается к /api на том же адресе, что и к страницам сайта.
import type { NextRequest } from "next/server";

import { backendURL } from "@/lib/server-api";

export const dynamic = "force-dynamic";

const maxBodySize = 16 * 1024;
const allowedRoutes = new Set([
  "health",
  "services",
  "csrf",
  "cookie-notice",
  "auth/register",
  "auth/login",
  "auth/logout",
  "auth/me",
  "requests",
]);

interface RouteContext {
  params: Promise<{ path: string[] }>;
}

class BodyTooLargeError extends Error {}

async function readLimitedBody(request: NextRequest): Promise<ArrayBuffer> {
  const reader = request.body?.getReader();

  if (!reader) {
    return new ArrayBuffer(0);
  }

  const chunks: Uint8Array[] = [];
  let size = 0;

  try {
    while (true) {
      const { done, value } = await reader.read();

      if (done) {
        break;
      }

      size += value.byteLength;

      if (size > maxBodySize) {
        await reader.cancel();
        throw new BodyTooLargeError();
      }

      chunks.push(value);
    }
  } finally {
    reader.releaseLock();
  }

  const body = new Uint8Array(size);
  let offset = 0;

  for (const chunk of chunks) {
    body.set(chunk, offset);
    offset += chunk.byteLength;
  }

  return body.buffer;
}

function forwardedHeaders(request: NextRequest): Headers {
  const headers = new Headers({ "Content-Type": "application/json" });

  for (const name of [
    "content-type",
    "cookie",
    "origin",
    "x-requested-with",
    "x-csrf-token",
  ]) {
    const value = request.headers.get(name);

    if (value) {
      headers.set(name, value);
    }
  }

  return headers;
}

function responseHeaders(response: Response): Headers {
  const headers = new Headers({
    "Content-Type": "application/json; charset=utf-8",
    "Cache-Control": "no-store",
  });

  for (const cookie of response.headers.getSetCookie()) {
    headers.append("Set-Cookie", cookie);
  }

  for (const name of ["retry-after", "content-security-policy"]) {
    const value = response.headers.get(name);

    if (value) {
      headers.set(name, value);
    }
  }

  return headers;
}

function errorResponse(message: string, status: number): Response {
  return Response.json(
    { error: message },
    { status, headers: { "Cache-Control": "no-store" } },
  );
}

async function forward(request: NextRequest, context: RouteContext) {
  const { path } = await context.params;
  const route = path.join("/");

  if (!allowedRoutes.has(route)) {
    return errorResponse("Страница не найдена.", 404);
  }

  try {
    const body = request.method === "POST"
      ? await readLimitedBody(request)
      : undefined;

    const response = await fetch(
      `${backendURL}/api/${route}${request.nextUrl.search}`,
      {
        method: request.method,
        cache: "no-store",
        headers: forwardedHeaders(request),
        body,
        signal: AbortSignal.timeout(10000),
      },
    );

    return new Response(response.body, {
      status: response.status,
      headers: responseHeaders(response),
    });
  } catch (error) {
    if (error instanceof BodyTooLargeError) {
      return errorResponse("Запрос превышает 16 КБ.", 413);
    }

    return errorResponse("Сервис временно недоступен. Попробуйте позже.", 503);
  }
}

export { forward as GET, forward as POST };
