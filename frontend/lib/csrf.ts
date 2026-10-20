// Подписанный double-submit токен; этот модуль используется только на сервере.
import { createHmac, randomBytes, timingSafeEqual } from "node:crypto";

export const csrfCookieName = "csrf_token";
export const csrfCookieOptions = {
  httpOnly: true,
  sameSite: "strict" as const,
  path: "/",
  maxAge: 7200,
  secure: process.env.COOKIE_SECURE === "true",
};

// signature отделяет назначение ключа от JWT и связывает токен с access-cookie.
function signature(nonce: string, session: string): string {
  const secret = process.env.JWT_SECRET || "";

  if (!/^[a-fA-F0-9]{64}$/.test(secret)) {
    throw new Error("JWT_SECRET должен содержать 32 байта в hex");
  }

  return createHmac("sha256", Buffer.from(secret, "hex"))
    .update(
      `studyblog:csrf:v1!${Buffer.byteLength(session)}!${session}!${nonce}`,
    )
    .digest("hex");
}

// createCSRFToken создаёт 256 случайных бит и подпись, не раскрывая JWT.
export function createCSRFToken(session = ""): string {
  const nonce = randomBytes(32).toString("hex");

  return `${nonce}.${signature(nonce, session)}`;
}

// validCSRFToken проверяет формат и подпись за постоянное время.
export function validCSRFToken(token: string, session = ""): boolean {
  if (!/^[a-f0-9]{64}\.[a-f0-9]{64}$/.test(token)) {
    return false;
  }

  const [nonce, mac] = token.split(".");

  return timingSafeEqual(
    Buffer.from(mac, "hex"),
    Buffer.from(signature(nonce, session), "hex"),
  );
}

// sameCSRFToken сравнивает cookie и поле формы без утечки совпадающего префикса.
export function sameCSRFToken(left: string, right: string): boolean {
  return (
    Buffer.byteLength(left) === 129 &&
    Buffer.byteLength(right) === 129 &&
    timingSafeEqual(Buffer.from(left), Buffer.from(right))
  );
}
