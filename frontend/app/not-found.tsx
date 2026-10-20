import Link from "next/link";
import { Sparkles } from "lucide-react";

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
