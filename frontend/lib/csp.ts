// Политика загрузки ресурсов для страниц и инструментов разработки Next.js.

/** Формирует CSP; разрешения DevTools включаются только для разработки. */
export function contentSecurityPolicy(
  nonce: string,
  development = false,
): string {
  const developmentScripts = development ? " 'unsafe-eval'" : "";

  return [
    "default-src 'self'",
    `script-src 'self' 'nonce-${nonce}'${developmentScripts}`,
    `style-src 'self' 'nonce-${nonce}'`,
    // DevTools создают теги style без nonce; атрибуты style остаются запрещены.
    ...(development ? ["style-src-elem 'self' 'unsafe-inline'"] : []),
    "img-src 'self' data:",
    "font-src 'self'",
    "connect-src 'self'",
    "object-src 'none'",
    "base-uri 'self'",
    "form-action 'self'",
    "frame-ancestors 'none'",
  ].join("; ");
}
