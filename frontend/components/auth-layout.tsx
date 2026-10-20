import Link from "next/link";
import { ArrowLeft, Check, Sparkles } from "lucide-react";

interface AuthLayoutProps {
  children: React.ReactNode;
  register: boolean;
}

export function AuthLayout({ children, register }: AuthLayoutProps) {
  const description = register
    ? "Заполните все поля. Для текстовых — до 20 символов, " +
      "для пароля ограничение в 20 символов не действует."
    : "Введите логин и пароль, чтобы запланировать уборку.";

  return (
    <section
      className={`container auth-layout ${register ? "registration-layout" : ""}`}
    >
      <aside className="auth-aside">
        <Link
          href="/"
          className="back-link"
        >
          <ArrowLeft size={16} /> На главную
        </Link>
        <span className="eyebrow">
          Чистый дом / {register ? "Знакомство" : "С возвращением"}
        </span>
        <h1>
          {register ? (
            <>
              Домашние дела
              <br />
              можно <em>делегировать.</em>
            </>
          ) : (
            <>
              Ваш чистый дом
              <br />
              <em>начинается здесь.</em>
            </>
          )}
        </h1>
        <p>
          Войдите в свой ритм. А для уборки найдётся время в нашем календаре.
        </p>
        <div className="auth-illustration">
          <Sparkles
            size={64}
            strokeWidth={1}
          />
          <span className="illustration-ring" />
          <span className="small-spark">
            <Sparkles size={27} />
          </span>
        </div>
        <ul className="auth-benefits">
          <li>
            <Check size={16} /> Удобная запись на уборку
          </li>
          <li>
            <Check size={16} /> Все заявки в одном месте
          </li>
          <li>
            <Check size={16} /> Только ваши данные и ваш кабинет
          </li>
        </ul>
      </aside>
      <div className="auth-panel">
        <span className="eyebrow">
          {register ? "Давайте знакомиться" : "Рады видеть вас снова"}
        </span>
        <h2>{register ? "Создать аккаунт" : "Вход в кабинет"}</h2>
        <p className="panel-description">{description}</p>
        {children}
      </div>
    </section>
  );
}
