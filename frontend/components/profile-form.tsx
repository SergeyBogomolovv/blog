"use client";

// Клиентская граница только для формы изменения профиля.
import { usePathname } from "next/navigation";
import { useActionState } from "react";

import { profileAction } from "@/app/actions";
import { Button } from "@/components/ui/button";

// ProfileForm сохраняет поля при отказе сервера и показывает результат действия.
export function ProfileForm({ name, bio }: { name: string; bio: string }) {
  const [state, action, pending] = useActionState(
    profileAction,
    {},
    usePathname(),
  );

  return (
    <form
      className="profile-form"
      action={action}
      autoComplete="off"
    >
      <h2>Немного о вас</h2>
      <label htmlFor="profile-name">Имя</label>
      <input
        id="profile-name"
        name="name"
        defaultValue={state.values?.name ?? name}
        required
        minLength={2}
        maxLength={60}
      />
      <label htmlFor="bio">О себе</label>
      <textarea
        id="bio"
        name="bio"
        defaultValue={state.values?.bio ?? bio}
        required
        minLength={1}
        maxLength={500}
        rows={3}
      />
      <span className="field-hint">
        До 500 символов. Имя и описание будут видны всем.
      </span>
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
      <Button
        type="submit"
        disabled={pending}
      >
        {pending ? "Сохраняем…" : "Сохранить изменения"}
      </Button>
    </form>
  );
}
