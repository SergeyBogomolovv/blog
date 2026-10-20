// Резервный документ для непредвиденных ошибок корневой разметки.
"use client";

import "./globals.css";
import { ServiceUnavailable } from "@/components/service-unavailable";

/** Сохраняет русскоязычный интерфейс даже при ошибке корневого layout. */
export default function GlobalError({ retry }: { retry: () => void }) {
  return (
    <html lang="ru">
      <head>
        <title>Сервис временно недоступен · Чистый дом</title>
      </head>
      <body>
        <ServiceUnavailable retry={retry} />
      </body>
    </html>
  );
}
