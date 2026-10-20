// Сообщение о временном сбое с повторной загрузкой текущей страницы.
"use client";

import Link from "next/link";

/** Параметры повторной загрузки после сбоя сервиса. */
interface ServiceUnavailableProps {
  retry?: () => void;
}

/** Перезагружает страницу, повторяя проверку сессии и получение данных. */
function reloadPage() {
  window.location.reload();
}

/** Показывает понятное сообщение, не раскрывая внутренние сведения сервера. */
export function ServiceUnavailable({
  retry = reloadPage,
}: ServiceUnavailableProps) {
  return (
    <section
      className="container error-page"
      aria-labelledby="unavailable-title"
    >
      <span className="eyebrow">Не удалось загрузить страницу</span>
      <h1 id="unavailable-title">Сервис временно недоступен</h1>
      <p role="status">
        Попробуйте обновить страницу через некоторое время. Ваши сохранённые
        заявки останутся в кабинете.
      </p>
      <button
        type="button"
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
    </section>
  );
}
