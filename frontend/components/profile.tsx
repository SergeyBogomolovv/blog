// Серверная публичная страница автора; форма видна только владельцу сессии.
import Link from "next/link";

import { Avatar } from "@/components/avatar";
import { Feed } from "@/components/blog";
import { ProfileForm } from "@/components/profile-form";
import { date, type User } from "@/lib/api";
import { getSession, readPublic } from "@/lib/server-api";

// Profile параллельно получает профиль и текущего пользователя без клиентского эффекта.
export async function Profile({ id, page }: { id: string; page: number }) {
  const [profile, user] = await Promise.all([
    readPublic<User>(`/users/${id}`),
    getSession(),
  ]);
  const own = user?.id === profile.id;

  return (
    <div className="profile-page">
      <section className="profile-header">
        <Avatar
          name={profile.name}
          large
        />
        <div>
          <span className="eyebrow">Страница автора</span>
          <h1>{profile.name}</h1>
          <p className="profile-login">
            @{profile.login}
            {profile.role === "admin" && (
              <span className="badge">Администратор</span>
            )}
            {profile.banned && (
              <span className="badge banned">Заблокирован</span>
            )}
          </p>
          <p className="profile-bio">{profile.bio}</p>
          <span className="field-hint">
            С нами с {date(profile.created_at)}
          </span>
        </div>
      </section>
      {own && (
        <details className="edit-profile">
          <summary>Редактировать профиль</summary>
          <ProfileForm
            name={profile.name}
            bio={profile.bio}
          />
        </details>
      )}
      <section className="profile-feed">
        <div className="section-heading">
          <h2>
            {own ? "Мои истории" : "Истории автора"}
            <span className="heading-dot">.</span>
          </h2>
          {own && (
            <Link
              href="/write"
              className="text-link"
            >
              Написать историю ↗
            </Link>
          )}
        </div>
        <Feed
          author={profile.id}
          page={page}
        />
      </section>
    </div>
  );
}
