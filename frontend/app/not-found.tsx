// Страница для неизвестных адресов приложения.
import Link from "next/link";

// NotFound предлагает вернуться в ленту после ошибочного перехода.
export default function NotFound() {
  return (
    <div className="empty-state">
      <span className="eyebrow">404</span>
      <h1>Эта строка пока не написана</h1>
      <p>Кажется, такой страницы здесь нет.</p>
      <Link
        href="/"
        className="write-link"
      >
        Вернуться к историям
      </Link>
    </div>
  );
}
