// Страница входа и подтверждение успешной регистрации.
import type { Metadata } from "next";
import { redirect } from "next/navigation";

import { AuthForm } from "@/components/auth-form";
import { AuthLayout } from "@/components/auth-layout";
import { getSession } from "@/lib/server-api";

export const metadata: Metadata = { title: "Вход" };

/** Параметры компонента LoginPageProps. */
interface LoginPageProps {
  searchParams: Promise<{ registered?: string; service?: string }>;
}

/** Показывает форму входа с учётом регистрации и выбранной услуги. */
export default async function Login({ searchParams }: LoginPageProps) {
  const parameters = await searchParams;
  const service =
    typeof parameters.service === "string" ? parameters.service : "";
  const serviceQuery = service ? `?service=${encodeURIComponent(service)}` : "";
  const user = await getSession();

  if (user) {
    redirect(`/dashboard${serviceQuery}`);
  }

  return (
    <AuthLayout register={false}>
      <AuthForm
        mode="login"
        registered={parameters.registered === "1"}
        service={service}
      />
    </AuthLayout>
  );
}
