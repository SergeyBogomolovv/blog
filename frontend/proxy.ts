// CSP и перенаправления выполняются до начала потокового рендера страниц.
import { randomBytes } from "node:crypto";
import { type NextRequest, NextResponse } from "next/server";

import { contentSecurityPolicy } from "@/lib/csp";
import {
  createCSRFToken,
  csrfCookieName,
  csrfCookieOptions,
  validCSRFToken,
} from "@/lib/csrf";

type SessionStatus = "guest" | "authenticated" | "unavailable";

interface SessionCheck {
  status: SessionStatus;
  cookies: string[];
}

async function checkSession(request: NextRequest): Promise<SessionCheck> {
  if (!request.cookies.has("cleaning.sid")) {
    return { status: "guest", cookies: [] };
  }

  if (request.method !== "GET" && request.method !== "HEAD") {
    return { status: "unavailable", cookies: [] };
  }

  const backendURL = process.env.BACKEND_URL || "http://127.0.0.1:8080";

  try {
    const response = await fetch(`${backendURL}/api/auth/me`, {
      headers: { Cookie: request.cookies.toString() },
      cache: "no-store",
      signal: AbortSignal.timeout(5000),
    });

    let status: SessionStatus = "unavailable";

    if (response.ok) {
      status = "authenticated";
    } else if (response.status === 401) {
      status = "guest";
    }

    return { status, cookies: response.headers.getSetCookie() };
  } catch {
    // При сбое API страница покажет ошибку; действующую cookie сохраняем.
    return { status: "unavailable", cookies: [] };
  }
}

function redirectDestination(
  request: NextRequest,
  session: SessionStatus,
): URL | null {
  const pathname = request.nextUrl.pathname;
  let destination: URL | null = null;

  if (pathname === "/dashboard" && session === "guest") {
    destination = new URL("/login", request.url);
  } else if (
    session === "authenticated" &&
    (pathname === "/login" || pathname === "/register")
  ) {
    destination = new URL("/dashboard", request.url);
  }

  const selectedService = request.nextUrl.searchParams.get("service");

  if (destination && selectedService) {
    destination.searchParams.set("service", selectedService);
  }

  return destination;
}

export async function proxy(request: NextRequest) {
  const nonce = randomBytes(18).toString("base64");
  const policy = contentSecurityPolicy(
    nonce,
    process.env.NODE_ENV === "development",
  );
  const sessionToken = request.cookies.get("cleaning.sid")?.value || "";
  let csrfToken = request.cookies.get(csrfCookieName)?.value || "";

  const isReadRequest = request.method === "GET" || request.method === "HEAD";
  const needsCSRFToken = isReadRequest && !validCSRFToken(csrfToken, sessionToken);

  if (needsCSRFToken) {
    csrfToken = createCSRFToken(sessionToken);
    request.cookies.set(csrfCookieName, csrfToken);
  }

  const headers = new Headers(request.headers);
  headers.set("Content-Security-Policy", policy);
  headers.set("x-nonce", nonce);

  const session = await checkSession(request);
  const destination = redirectDestination(request, session.status);
  const response = destination
    ? NextResponse.redirect(destination)
    : NextResponse.next({ request: { headers } });

  response.headers.set("Content-Security-Policy", policy);
  response.headers.set("Cache-Control", "no-store");

  if (needsCSRFToken) {
    response.cookies.set(csrfCookieName, csrfToken, csrfCookieOptions);
  }

  for (const cookie of session.cookies) {
    response.headers.append("Set-Cookie", cookie);
  }

  return response;
}

export const config = {
  matcher: [
    "/((?!api|_next/static|_next/image|favicon.ico|icon.svg|clean-room.svg).*)",
  ],
};
