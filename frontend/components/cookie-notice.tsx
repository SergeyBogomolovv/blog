"use client";

// Только управление нативным диалогом; решение о показе принимает сервер по cookie.
import { usePathname } from "next/navigation";
import { useActionState } from "react";

import { acceptCookiesAction } from "@/app/actions";
import { Button } from "@/components/ui/button";
import { CSRFField } from "@/components/csrf";

// openDialog открывает модальное окно после присоединения DOM-элемента.
function openDialog(node: HTMLDialogElement | null) {
  if (node && !node.open) {
    node.showModal();
  }
}

// CookieNotice отправляет ознакомление серверным действием без localStorage и эффектов.
export function CookieNotice() {
  const [state, action, pending] = useActionState(
    acceptCookiesAction,
    {},
    usePathname(),
  );

  return (
    <dialog
      ref={openDialog}
      className="cookie-dialog"
      aria-labelledby="cookie-title"
      onCancel={(event) => event.preventDefault()}
    >
      <span className="eyebrow">Небольшая формальность</span>
      <h2 id="cookie-title">Здесь есть cookie</h2>
      <p>
        Мы используем только необходимые cookie, чтобы сохранять ваш вход и
        отметку об ознакомлении, а также защищать формы от поддельных запросов.
        Рекламных и аналитических cookie нет.
      </p>
      <form
        action={action}
        autoComplete="off"
      >
        <CSRFField />
        {state.error && (
          <p
            className="error-box"
            role="alert"
          >
            {state.error}
          </p>
        )}
        <button
          type="submit"
          name="destination"
          value="privacy"
          className="text-link"
          disabled={pending}
        >
          Политика обработки данных →
        </button>
        <Button
          type="submit"
          className="wide-button"
          disabled={pending}
        >
          {pending ? "Подождите…" : "Понятно, продолжить"}
        </Button>
      </form>
    </dialog>
  );
}
