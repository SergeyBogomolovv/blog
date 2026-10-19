// Серверный корневой макет с русским языком и метаданными.
import "./globals.css";

import type { Metadata } from "next";

import { Shell } from "@/components/shell";
export const metadata: Metadata = {
  title: {
    default: "Между строк — место для ваших историй",
    template: "%s · Между строк",
  },
  description: "Небольшой блог о мыслях, открытиях и повседневных историях.",
};

// RootLayout задаёт общую структуру всех страниц.
export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html
      lang="ru"
      data-scroll-behavior="smooth"
    >
      <body>
        <Shell>{children}</Shell>
      </body>
    </html>
  );
}
