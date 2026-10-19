// Серверная лента и карточки публикаций: данные включены в HTML до гидратации.
import Link from "next/link";

import { ArrowRight, BookOpen, MessageCircle, PenLine } from "lucide-react";

import { Avatar } from "@/components/avatar";
import { date, type Post } from "@/lib/api";
import { getSession, readPublic } from "@/lib/server-api";

// PostCard представляет публикацию в общей ленте и на странице пользователя.
export function PostCard({ post }: { post: Post }) {
  return (
    <article className="post-card">
      <div className="post-meta">
        <Link
          className="author-link"
          href={`/users/${post.author.id}`}
        >
          <Avatar name={post.author.name} />
          <span>{post.author.name}</span>
        </Link>
        <span className="meta-dot">·</span>
        <time dateTime={post.created_at}>{date(post.created_at)}</time>
      </div>
      <Link
        href={`/posts/${post.id}`}
        className="post-title-link"
      >
        <h2>{post.title}</h2>
      </Link>
      <p className="post-excerpt">{post.body}</p>
      <div className="post-bottom">
        <Link
          href={`/posts/${post.id}`}
          className="read-link"
        >
          Читать историю
          <ArrowRight size={15} />
        </Link>
        <Link
          href={`/posts/${post.id}#comments`}
          className="comment-count"
          aria-label={`Комментарии: ${post.comment_count}`}
        >
          <MessageCircle size={16} />
          {post.comment_count}
        </Link>
      </div>
    </article>
  );
}

// Feed загружает страницу публикаций на сервере; ссылки работают без клиентского состояния.
export async function Feed({
  author,
  page = 0,
}: {
  author?: number;
  page?: number;
}) {
  const posts = await readPublic<Post[]>(
    `/posts?page=${page}${author ? `&author=${author}` : ""}`,
  );
  const base = author ? `/users/${author}` : "/";

  return (
    <>
      {posts.map((post) => (
        <PostCard
          key={post.id}
          post={post}
        />
      ))}
      {posts.length === 0 && (
        <div className="empty-state">
          <span className="empty-icon">
            <BookOpen
              size={30}
              strokeWidth={1.4}
            />
          </span>
          <h2>
            {page ? "Здесь больше нет историй" : "У каждой истории есть начало"}
          </h2>
          <p>Здесь пока нет публикаций.</p>
          <Link
            className="text-link"
            href={page ? base : "/write"}
          >
            {page ? "К первой странице" : "Написать первую историю"}
            <ArrowRight size={16} />
          </Link>
        </div>
      )}
      <nav
        className="pagination"
        aria-label="Страницы публикаций"
      >
        {page > 0 && (
          <Link
            className="text-link"
            href={`${base}?page=${page - 1}`}
          >
            ← Предыдущая
          </Link>
        )}
        {(page > 0 || posts.length === 20) && <span>Страница {page + 1}</span>}
        {posts.length === 20 && (
          <Link
            className="text-link"
            href={`${base}?page=${page + 1}`}
          >
            Следующая →
          </Link>
        )}
      </nav>
    </>
  );
}

// HomePage объединяет ленту и короткое приглашение к участию.
export async function HomePage({ page }: { page: number }) {
  const user = await getSession();

  return (
    <>
      <section className="hero">
        <div>
          <div className="eyebrow">
            <span className="small-dot" />
            Обычные люди. Необычные мысли.
          </div>
          <h1>
            Истории, которые
            <br />
            остаются <em>между строк.</em>
          </h1>
          <p>
            Открывайте новое, делитесь своим и находите
            <br className="desktop-break" /> тех, кто смотрит на мир немного
            похоже.
          </p>
        </div>
        <div
          className="hero-art"
          aria-hidden="true"
        >
          <div className="art-line" />
          <span className="art-quote">“</span>
          <div className="art-note">
            у каждой мысли
            <br />
            есть своя история<span>и место здесь ↗</span>
          </div>
          <span className="art-star">✳</span>
        </div>
      </section>
      <div className="content-grid">
        <section aria-labelledby="feed-title">
          <div className="section-heading">
            <h2 id="feed-title">
              Свежие истории<span className="heading-dot">.</span>
            </h2>
          </div>
          <Feed page={page} />
        </section>
        <aside className="sidebar">
          <div className="welcome-card">
            <span className="card-icon">
              <PenLine size={23} />
            </span>
            <h2>
              Есть чем
              <br />
              поделиться?
            </h2>
            <p>
              Наблюдение, маленькое открытие или просто мысль вслух. Здесь всему
              найдётся место.
            </p>
            <Link
              className="write-link"
              href={user ? "/write" : "/register"}
            >
              <span>{user ? "Написать историю" : "Стать автором"}</span>
              <ArrowRight size={17} />
            </Link>
          </div>
          <div className="about-card">
            <span className="eyebrow">О нашем пространстве</span>
            <p>
              Без гонки за лайками.
              <br />С интересом друг к другу.
            </p>
            <ul>
              <li>Пишите о том, что важно вам</li>
              <li>Обсуждайте, уважая собеседника</li>
              <li>Оставайтесь собой</li>
            </ul>
          </div>
        </aside>
      </div>
    </>
  );
}
