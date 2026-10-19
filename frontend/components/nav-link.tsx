"use client";

// Клиентская граница только для подсветки активного маршрута.
import Link from "next/link";
import { usePathname } from "next/navigation";
import type { ReactNode } from "react";

import { cn } from "@/lib/utils";

// NavLink определяет активность по текущему URL без собственного состояния.
export function NavLink({
  href,
  children,
  className,
}: {
  href: string;
  children: ReactNode;
  className?: string;
}) {
  const active = usePathname() === href;

  return (
    <Link
      href={href}
      className={cn("nav-link", active && "active", className)}
      aria-current={active ? "page" : undefined}
    >
      {children}
    </Link>
  );
}
