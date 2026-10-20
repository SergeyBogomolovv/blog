// Страница для отсутствующих адресов сайта.
import Link from "next/link";
import { Sparkles } from "lucide-react";

/** Показывает сообщение об отсутствующей странице и ссылку на главную. */
export default function NotFound() {
  return (
    <div className="container error-page">
      <Sparkles size={40} />
      <span className="eyebrow">404 / Страница не найдена</span>
      <h1>
        Здесь пока <em>пусто.</em>
      </h1>
      <p>Вернитесь на главную, чтобы выбрать уборку или открыть кабинет.</p>
      <Link
        className="button"
        href="/"
      >
        На главную
      </Link>
    </div>
  );
}
