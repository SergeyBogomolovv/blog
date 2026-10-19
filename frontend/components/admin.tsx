// Серверная панель модерации; список аккаунтов не загружается клиентским эффектом.
import Link from "next/link";
import { redirect } from "next/navigation";

import { banAction } from "@/app/actions";
import { ActionForm } from "@/components/action-form";
import { Avatar } from "@/components/avatar";
import type { User } from "@/lib/api";
import { getSession, readAPI } from "@/lib/server-api";

// Admin проверяет роль до загрузки списка; Server Action и Go проверяют её повторно.
export async function Admin() {
  const user = await getSession();

  if (!user) {
    redirect("/login");
  }

  if (user.role !== "admin") {
    return (
      <div className="empty-state">
        <h1>Доступ только для администратора</h1>
        <Link
          href="/"
          className="text-link"
        >
          Вернуться к историям
        </Link>
      </div>
    );
  }

  const users = await readAPI<User[]>("/admin/users");

  return (
    <section className="admin-page">
      <span className="eyebrow">Забота о пространстве</span>
      <h1>Модерация</h1>
      <p className="form-intro">
        Блокируйте нарушителей. Удаление публикации доступно на странице самой
        истории.
      </p>
      <div className="admin-list">
        {users.map((account) => (
          <div
            key={account.id}
            className="admin-row"
          >
            <Link
              className="author-link"
              href={`/users/${account.id}`}
            >
              <Avatar name={account.name} />
              <span>
                <strong>{account.name}</strong>
                <small>@{account.login}</small>
              </span>
            </Link>
            <span className={`badge${account.banned ? " banned" : ""}`}>
              {account.role === "admin"
                ? "Администратор"
                : account.banned
                  ? "Заблокирован"
                  : "Участник"}
            </span>
            {account.role !== "admin" && (
              <details
                key={`${account.id}-${account.banned}`}
                className="moderation-action"
              >
                <summary>
                  {account.banned ? "Разблокировать" : "Заблокировать"}
                </summary>
                <ActionForm
                  action={banAction}
                  label="Подтвердить"
                  variant="outline"
                >
                  <input
                    type="hidden"
                    name="user_id"
                    value={account.id}
                  />
                  <input
                    type="hidden"
                    name="banned"
                    value={String(!account.banned)}
                  />
                  <label className="checkbox-label">
                    <input
                      type="checkbox"
                      name="confirm"
                      required
                    />
                    <span>
                      {account.banned
                        ? "Разрешить вход в аккаунт"
                        : "Заблокировать аккаунт и завершить активные сессии"}
                    </span>
                  </label>
                </ActionForm>
              </details>
            )}
          </div>
        ))}
      </div>
      <p className="field-hint">Показаны последние 200 аккаунтов.</p>
    </section>
  );
}
