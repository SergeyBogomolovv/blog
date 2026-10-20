"use client";

// Форма комментария использует результат Server Action вместо отдельной загрузки списка.
import { useActionState } from "react";

import { commentAction } from "@/app/actions";
import { Button } from "@/components/ui/button";
import { CSRFField } from "@/components/csrf";

// CommentForm сохраняет неотправленный текст при ошибке и очищает поле после публикации.
export function CommentForm({ postID }: { postID: number }) {
  const [state, action, pending] = useActionState(
    commentAction,
    {},
    `/posts/${postID}`,
  );

  return (
    <form
      className="comment-form"
      action={action}
      autoComplete="off"
    >
      <CSRFField />
      <input
        type="hidden"
        name="post_id"
        value={postID}
      />
      <label htmlFor="comment">Поделитесь мыслью</label>
      <textarea
        id="comment"
        name="body"
        defaultValue={state.values?.body ?? ""}
        required
        minLength={1}
        maxLength={2000}
        rows={3}
        placeholder="Что откликнулось в этой истории?"
      />
      {state.error && (
        <p
          className="error-box"
          role="alert"
        >
          {state.error}
        </p>
      )}
      {state.success && (
        <p
          className="success-box"
          role="status"
        >
          {state.success}
        </p>
      )}
      <div className="comment-submit">
        <span className="field-hint">
          До 2 000 символов. Будьте доброжелательны.
        </span>
        <Button
          type="submit"
          disabled={pending}
        >
          {pending ? "Отправляем…" : "Отправить"}
        </Button>
      </div>
    </form>
  );
}
