"use client";

// Редактор отправляет публикацию через Server Action без клиентского HTTP-клиента.
import Link from "next/link";
import { useActionState } from "react";

import { ArrowLeft, ArrowUpRight } from "lucide-react";

import { createPostAction } from "@/app/actions";
import { Button } from "@/components/ui/button";
import { CSRFField } from "@/components/csrf";

// Editor сохраняет введённый текст при ошибке и показывает состояние отправки.
export function Editor() {
  const [state, formAction, busy] = useActionState(
    createPostAction,
    {},
    "/write",
  );

  return (
    <section className="editor-page">
      <Link
        className="back-link"
        href="/"
      >
        <ArrowLeft size={16} />К историям
      </Link>
      <span className="eyebrow">Чистый лист — хорошее начало</span>
      <h1>О чём ваша история?</h1>
      <p className="form-intro">
        Пишите свободно. Иногда самое простое — самое интересное.
      </p>
      <form
        action={formAction}
        autoComplete="off"
      >
        <CSRFField />
        <label htmlFor="title">Заголовок</label>
        <input
          id="title"
          name="title"
          defaultValue={state.values?.title ?? ""}
          required
          minLength={3}
          maxLength={120}
          placeholder="Дайте вашей мысли название"
          aria-describedby="title-hint"
        />
        <span
          className="field-hint"
          id="title-hint"
        >
          От 3 до 120 символов.
        </span>
        <label htmlFor="body">Ваша история</label>
        <textarea
          id="body"
          name="body"
          defaultValue={state.values?.body ?? ""}
          required
          minLength={1}
          maxLength={10000}
          rows={14}
          placeholder="Всё началось с…"
        />
        <span className="field-hint">
          До 10 000 символов. Текст будет виден всем посетителям.
        </span>
        {state.error && (
          <p
            className="error-box"
            role="alert"
          >
            {state.error}
          </p>
        )}
        <div className="editor-bottom">
          <span>Ваши слова имеют значение.</span>
          <Button
            type="submit"
            disabled={busy}
          >
            {busy ? "Публикуем…" : "Опубликовать"}
            <ArrowUpRight size={17} />
          </Button>
        </div>
      </form>
    </section>
  );
}
