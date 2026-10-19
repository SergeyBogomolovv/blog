"use client";

// Мобильное меню использует нативное раскрытие без хранения серверных данных на клиенте.
import { usePathname } from "next/navigation";
import {
  type KeyboardEvent,
  type MouseEvent,
  type ReactNode,
  useRef,
} from "react";

import { Menu, X } from "lucide-react";

// MobileMenu закрывает панель при переходе по ссылке, смене маршрута или нажатии Escape.
export function MobileMenu({ children }: { children: ReactNode }) {
  const pathname = usePathname();
  const menu = useRef<HTMLDetailsElement>(null);

  // closeOnNavigation закрывает панель, сохраняя обычное поведение ссылки Next.js.
  function closeOnNavigation(event: MouseEvent<HTMLDivElement>) {
    if (
      event.target instanceof Element &&
      event.target.closest("a") &&
      menu.current
    ) {
      menu.current.open = false;
    }
  }

  // closeOnEscape возвращает клавиатурный фокус на кнопку меню.
  function closeOnEscape(event: KeyboardEvent<HTMLDetailsElement>) {
    if (event.key === "Escape" && menu.current?.open) {
      event.preventDefault();
      menu.current.open = false;
      menu.current.querySelector("summary")?.focus();
    }
  }

  return (
    <details
      key={pathname}
      ref={menu}
      className="mobile-menu"
      onKeyDown={closeOnEscape}
    >
      <summary
        className="mobile-menu-toggle"
        aria-label="Меню навигации"
      >
        <Menu
          className="menu-open-icon"
          size={22}
          aria-hidden="true"
        />
        <X
          className="menu-close-icon"
          size={22}
          aria-hidden="true"
        />
      </summary>
      <div
        className="mobile-menu-panel"
        onClick={closeOnNavigation}
      >
        {children}
      </div>
    </details>
  );
}
