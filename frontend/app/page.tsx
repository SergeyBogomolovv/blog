// Лента рендерится на сервере; номер страницы хранится в URL.
import { Suspense } from "react";

import { HomePage } from "@/components/blog";
import { pageNumber } from "@/lib/params";

// Page передаёт проверенный параметр пагинации серверному компоненту.
export default async function Page({
  searchParams,
}: {
  searchParams: Promise<{ page?: string | string[] }>;
}) {
  const page = pageNumber((await searchParams).page);

  return (
    <Suspense
      key={page}
      fallback={
        <p
          role="status"
          className="loading"
        >
          Собираем истории…
        </p>
      }
    >
      <HomePage page={page} />
    </Suspense>
  );
}
