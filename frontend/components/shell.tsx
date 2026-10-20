// Общие навигация, основная область и подвал сайта.
import Link from "next/link";
import { ArrowUpRight, Sparkles } from "lucide-react";

import { LogoutButton } from "@/components/logout-button";
import type { User } from "@/lib/api";

/** Параметры компонента ShellProps. */
interface ShellProps {
  children: React.ReactNode;
  user: User | null;
}

/** Отображает общую оболочку с навигацией для гостя или пользователя. */
export function Shell({ children, user }: ShellProps) {
  return (
    <>
      <a
        className="skip-link"
        href="#main"
      >
        Перейти к содержимому
      </a>
      <header className="site-header">
        <div className="container header-inner">
          <Link
            href="/"
            className="brand"
            aria-label="Чистый дом — главная"
          >
            <span className="brand-mark">
              <Sparkles size={23} />
            </span>
            <span>
              Чистый дом
              <span className="brand-caption">Забота о вашем пространстве</span>
            </span>
          </Link>
          <nav aria-label="Основная навигация">
            <Link
              className="nav-item"
              href="/#services"
            >
              Услуги
            </Link>
            <Link
              className="nav-item"
              href="/#how-it-works"
            >
              Как это работает
            </Link>
            {user ? (
              <>
                <Link
                  className="button button-small"
                  href="/dashboard"
                >
                  Личный кабинет <ArrowUpRight size={16} />
                </Link>
                <LogoutButton />
              </>
            ) : (
              <>
                <Link
                  className="login-link"
                  href="/login"
                >
                  Войти
                </Link>
                <Link
                  className="button button-small header-cta"
                  href="/register"
                >
                  Регистрация <ArrowUpRight size={16} />
                </Link>
              </>
            )}
          </nav>
        </div>
      </header>
      <main id="main">{children}</main>
      <footer className="site-footer">
        <div className="container footer-inner">
          <Link
            className="brand footer-brand"
            href="/"
          >
            <Sparkles size={24} />
            Чистый дом
          </Link>
          <p>Меньше домашних дел. Больше жизни.</p>
          <Link href="/privacy">Обработка персональных данных</Link>
        </div>
        <div className="container footer-bottom">
          <span>© 2026 Чистый дом</span>
          <span>Учебный проект · Богомолов Сергей · ИСП-406</span>
        </div>
      </footer>
    </>
  );
}
