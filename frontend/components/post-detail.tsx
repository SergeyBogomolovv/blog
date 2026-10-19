// Серверная публикация и комментарии; клиентскими остаются только формы действий.
import Link from "next/link";

import { ArrowLeft, MessageCircle } from "lucide-react";

import { deletePostAction } from "@/app/actions";
import { ActionForm } from "@/components/action-form";
import { Avatar } from "@/components/avatar";
import { CommentForm } from "@/components/comment-form";
import { type Comment, date, type Post } from "@/lib/api";
import { getSession, readPublic } from "@/lib/server-api";

// PostDetail получает пост, обсуждение и сессию параллельно на сервере.
export async function PostDetail({ id }: { id: string }) {
  const [post, comments, user] = await Promise.all([
    readPublic<Post>(`/posts/${id}`),
    readPublic<Comment[]>(`/posts/${id}/comments`),
    getSession(),
  ]);

  return (
    <div className="article-page">
      <Link
        className="back-link"
        href="/"
      >
        <ArrowLeft size={16} />
        Все истории
      </Link>
      <article>
        <div className="post-meta">
          <Link
            className="author-link"
            href={`/users/${post.author.id}`}
          >
            <Avatar name={post.author.name} />
            <span>{post.author.name}</span>
          </Link>
          <span>·</span>
          <time dateTime={post.created_at}>{date(post.created_at)}</time>
        </div>
        <h1>{post.title}</h1>
        <div className="article-body">{post.body}</div>
        <div className="article-signature">
          История от{" "}
          <Link href={`/users/${post.author.id}`}>@{post.author.login}</Link>
        </div>
      </article>
      {(user?.id === post.author_id || user?.role === "admin") && (
        <details className="delete-area">
          <summary>Удалить историю</summary>
          <ActionForm
            action={deletePostAction}
            label="Удалить историю"
            variant="destructive"
          >
            <input
              type="hidden"
              name="post_id"
              value={post.id}
            />
            <label className="checkbox-label">
              <input
                type="checkbox"
                name="confirm"
                required
              />
              <span>
                Удалить историю и все комментарии без возможности восстановления
              </span>
            </label>
          </ActionForm>
        </details>
      )}
      <section
        id="comments"
        className="comments-section"
      >
        <h2>
          <MessageCircle size={23} />
          Обсуждение <span>{comments.length}</span>
        </h2>
        {user ? (
          <CommentForm postID={post.id} />
        ) : (
          <div className="guest-notice">
            <Link href="/login">Войдите</Link> или{" "}
            <Link href="/register">создайте аккаунт</Link>, чтобы участвовать в
            обсуждении.
          </div>
        )}
        {comments.length === 0 && (
          <p className="no-comments">Здесь пока тихо. Начните разговор.</p>
        )}
        {comments.map((comment) => (
          <article
            key={comment.id}
            className="comment"
          >
            <Link
              className="comment-avatar"
              href={`/users/${comment.author.id}`}
              aria-label={`Профиль ${comment.author.name}`}
            >
              <Avatar name={comment.author.name} />
            </Link>
            <div className="comment-content">
              <div className="post-meta">
                <Link
                  className="author-link"
                  href={`/users/${comment.author.id}`}
                >
                  {comment.author.name}
                </Link>
                <time dateTime={comment.created_at}>
                  {date(comment.created_at)}
                </time>
              </div>
              <p>{comment.body}</p>
            </div>
          </article>
        ))}
      </section>
    </div>
  );
}
