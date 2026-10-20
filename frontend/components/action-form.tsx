"use client";

// Небольшая клиентская граница для состояния отправки Server Action и его результата.
import { usePathname } from "next/navigation";
import { type ReactNode, useActionState } from "react";

import { Button } from "@/components/ui/button";
import { CSRFField } from "@/components/csrf";
import type { FormState } from "@/lib/form-state";

// ActionForm связывает серверное действие с React 19 без ручного fetch и синхронизации состояния.
export function ActionForm({
  action,
  children,
  label,
  className,
  buttonClassName,
  variant = "default",
}: {
  action: (state: FormState, data: FormData) => Promise<FormState>;
  children?: ReactNode;
  label: string;
  className?: string;
  buttonClassName?: string;
  variant?: "default" | "outline" | "destructive" | "ghost";
}) {
  const [state, formAction, pending] = useActionState(
    action,
    {},
    usePathname(),
  );

  return (
    <form
      action={formAction}
      className={className}
      autoComplete="off"
    >
      <CSRFField />
      {children}
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
        className={buttonClassName}
        variant={variant}
        disabled={pending}
      >
        {pending ? "Подождите…" : label}
      </Button>
    </form>
  );
}
