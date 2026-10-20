// Резервное отображение непредвиденных ошибок внутри страниц.
"use client";

import { ServiceUnavailable } from "@/components/service-unavailable";

/** Позволяет повторить получение данных после ошибки страницы. */
export default function ErrorPage({ retry }: { retry: () => void }) {
  return <ServiceUnavailable retry={retry} />;
}
