// Серверная навигация получает текущего пользователя из проверенной cookie.
import { cookies } from "next/headers";
import Link from "next/link";
import type { ReactNode } from "react";

import { ArrowUpRight, Feather, PenLine, Shield } from "lucide-react";

import { logoutAction } from "@/app/actions";
import { ActionForm } from "@/components/action-form";
import { CookieNotice } from "@/components/cookie-notice";
import { MobileMenu } from "@/components/mobile-menu";
import { NavLink } from "@/components/nav-link";
import type { User } from "@/lib/api";
import { getSession } from "@/lib/server-api";

// HeaderLinks формирует единый набор ссылок для обычного и мобильного меню.
function HeaderLinks({ user }: { user: User | null }) {
  return (
    <>
      <nav
        className="header-nav"
        aria-label="Главная навигация"
      >
        <NavLink href="/">Лента</NavLink>
        {user && <NavLink href={`/users/${user.id}`}>Мой профиль</NavLink>}
        {user?.role === "admin" && (
          <NavLink href="/admin">
            <Shield size={16} />
            Модерация
          </NavLink>
        )}
      </nav>
      <div className="header-actions">
        {user ? (
          <>
            <Link
              className="write-link header-action"
              href="/write"
            >
              <PenLine size={16} />
              <span>Написать</span>
            </Link>
            <ActionForm
              action={logoutAction}
              label="Выйти"
              variant="ghost"
              buttonClassName="header-action"
            />
          </>
        ) : (
          <>
            <NavLink
              href="/login"
              className="header-action"
            >
              Войти
            </NavLink>
            <Link
              className="write-link header-action"
              href="/register"
            >
              <span>Присоединиться</span>
              <ArrowUpRight
                size={18}
                aria-hidden="true"
              />
            </Link>
          </>
        )}
      </div>
    </>
  );
}

// Shell формирует общий каркас; персональные данные не кешируются между пользователями.
export async function Shell({ children }: { children: ReactNode }) {
  const [user, jar] = await Promise.all([
    getSession().catch(() => null),
    cookies(),
  ]);

  return (
    <>
      <a
        className="skip-link"
        href="#content"
      >
        Перейти к содержимому
      </a>
      <header className="site-header">
        <div className="header-inner">
          <Link
            href="/"
            className="brand"
            aria-label="Между строк — главная"
          >
            <span className="brand-mark">
              <Feather size={23} />
            </span>
            <span>
              между строк<span className="brand-dot">.</span>
            </span>
          </Link>
          <div className="header-desktop">
            <HeaderLinks user={user} />
          </div>
          <MobileMenu key={user?.id ?? "guest"}>
            <HeaderLinks user={user} />
          </MobileMenu>
        </div>
      </header>
      <main
        id="content"
        className="site-main"
      >
        {children}
      </main>
      <footer className="site-footer">
        <Link
          className="footer-brand"
          href="/"
        >
          между строк.
        </Link>
        <span>Место для мыслей, которыми хочется поделиться.</span>
        <Link href="/privacy">Данные и cookie</Link>
      </footer>
      {jar.get("cookie_notice")?.value !== "accepted" && <CookieNotice />}
    </>
  );
}
