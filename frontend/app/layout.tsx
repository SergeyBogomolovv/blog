// Общая разметка, метаданные и навигация страниц сайта.
import "./globals.css";

import type { Metadata } from "next";
import { cookies } from "next/headers";

import { Shell } from "@/components/shell";
import { CookieNotice } from "@/components/cookie-notice";
import { ServiceUnavailable } from "@/components/service-unavailable";
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

/** Параметры компонента RootLayoutProps. */
interface RootLayoutProps {
  children: React.ReactNode;
}

/** Добавляет общую оболочку с учётом авторизации пользователя. */
export default async function RootLayout({ children }: RootLayoutProps) {
  const session = await getSession();
  const showCookieNotice = !(await cookies()).has("cookie_notice");

  return (
    <html
      lang="ru"
      data-scroll-behavior="smooth"
    >
      <body>
        <Shell user={session.user}>
          {session.unavailable ? <ServiceUnavailable /> : children}
        </Shell>
        {showCookieNotice && <CookieNotice />}
      </body>
    </html>
  );
}
