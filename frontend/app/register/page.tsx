// Страница регистрации с сохранением выбранной услуги.
import type { Metadata } from "next";
import { redirect } from "next/navigation";

import { AuthForm } from "@/components/auth-form";
import { AuthLayout } from "@/components/auth-layout";
import { getSession } from "@/lib/server-api";

export const metadata: Metadata = { title: "Регистрация" };

/** Параметры компонента RegisterPageProps. */
interface RegisterPageProps {
  searchParams: Promise<{ service?: string }>;
}

/** Передаёт выбранную услугу в форму регистрации. */
export default async function Register({ searchParams }: RegisterPageProps) {
  const parameters = await searchParams;
  const service =
    typeof parameters.service === "string" ? parameters.service : "";
  const serviceQuery = service ? `?service=${encodeURIComponent(service)}` : "";
  const user = await getSession();

  if (user) {
    redirect(`/dashboard${serviceQuery}`);
  }

  return (
    <AuthLayout register>
      <AuthForm
        mode="register"
        service={service}
      />
    </AuthLayout>
  );
}
