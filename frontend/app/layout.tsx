// Общая разметка, метаданные и навигация страниц сайта.
import "./globals.css";

import type { Metadata } from "next";

import { Shell } from "@/components/shell";
import { getSession } from "@/lib/server-api";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: {
    default: "Чистый дом — забота о вашем пространстве",
    template: "%s · Чистый дом",
  },
  description:
    "Поддерживающая, генеральная уборка и уборка после ремонта. Выберите услугу и оставьте заявку в личном кабинете.",
};

interface RootLayoutProps {
  children: React.ReactNode;
}

/** Добавляет общую оболочку с учётом авторизации пользователя. */
export default async function RootLayout({ children }: RootLayoutProps) {
  const user = await getSession();

  return (
    <html lang="ru" data-scroll-behavior="smooth">
      <body>
        <Shell user={user}>{children}</Shell>
      </body>
    </html>
  );
}
