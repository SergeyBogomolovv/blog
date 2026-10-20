"use client";

import Link from "next/link";

export default function ErrorPage({ retry }: { retry: () => void }) {
  return (
    <div className="container error-page">
      <span className="eyebrow">Не удалось загрузить страницу</span>
      <h1>
        Попробуем <em>ещё раз?</em>
      </h1>
      <p>Проверьте, что Go API и база данных запущены.</p>
      <button
        className="button"
        onClick={retry}
      >
        Повторить загрузку
      </button>
      <Link
        className="text-link"
        href="/"
      >
        На главную
      </Link>
    </div>
  );
}
