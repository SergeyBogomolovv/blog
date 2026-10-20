// Выдача CSRF-cookie до рендера, включая первое открытие формы без JavaScript.
import { type NextRequest, NextResponse } from "next/server";

import {
  createCSRFToken,
  csrfCookieName,
  csrfCookieOptions,
  validCSRFToken,
} from "@/lib/csrf";

export function proxy(request: NextRequest) {
  // POST не получает новый токен: действие обязано проверить присланный браузером.
  if (request.method !== "GET" && request.method !== "HEAD") {
    return NextResponse.next();
  }

  const session = request.cookies.get("access_token")?.value || "";
  const token = request.cookies.get(csrfCookieName)?.value || "";

  if (validCSRFToken(token, session)) {
    return NextResponse.next();
  }

  const nextToken = createCSRFToken(session);
  request.cookies.set(csrfCookieName, nextToken);
  const response = NextResponse.next({ request: { headers: request.headers } });
  response.cookies.set(csrfCookieName, nextToken, csrfCookieOptions);
  response.headers.set("Cache-Control", "private, no-store");

  return response;
}

export const config = {
  matcher: ["/((?!_next/static|_next/image|favicon.ico|icon.svg).*)"],
};
