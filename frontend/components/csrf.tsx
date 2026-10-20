"use client";

// В клиент передаётся только CSRF-токен; access-cookie остаётся на сервере.
import { createContext, type ReactNode, useContext } from "react";

const CSRFContext = createContext("");

export function CSRFProvider({
  token,
  children,
}: {
  token: string;
  children: ReactNode;
}) {
  return <CSRFContext.Provider value={token}>{children}</CSRFContext.Provider>;
}

// CSRFField включается во все формы, вызывающие серверные действия.
export function CSRFField() {
  return (
    <input
      type="hidden"
      name="csrf_token"
      value={useContext(CSRFContext)}
    />
  );
}
