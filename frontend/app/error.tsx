"use client";

// Граница непредвиденных ошибок серверного рендера без раскрытия внутренних деталей.
import Link from "next/link";

import { Button } from "@/components/ui/button";

// ErrorPage позволяет повторить рендер после восстановления связи с API.
export default function ErrorPage({ retry }: { retry: () => void }) {
  return (
    <div className="empty-state">
      <h1>Не удалось открыть страницу</h1>
      <p>Проверьте, что бэкенд запущен, и попробуйте ещё раз.</p>
      <Button onClick={retry}>Попробовать снова</Button>
      <p>
        <Link
          href="/"
          className="text-link"
        >
          К историям
        </Link>
      </p>
    </div>
  );
}
