// Серверный корневой макет с русским языком и метаданными.
import "./globals.css";

import type { Metadata } from "next";
import { cookies } from "next/headers";

import { CSRFProvider } from "@/components/csrf";
import { Shell } from "@/components/shell";
import { csrfCookieName } from "@/lib/csrf";
export const metadata: Metadata = {
  title: {
    default: "Между строк — место для ваших историй",
    template: "%s · Между строк",
  },
  description: "Небольшой блог о мыслях, открытиях и повседневных историях.",
};

// RootLayout задаёт общую структуру всех страниц.
export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const token = (await cookies()).get(csrfCookieName)?.value || "";

  return (
    <html
      lang="ru"
      data-scroll-behavior="smooth"
    >
      <body>
        <CSRFProvider token={token}>
          <Shell>{children}</Shell>
        </CSRFProvider>
      </body>
    </html>
  );
}
