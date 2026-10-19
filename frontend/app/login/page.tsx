// Страница: Вход.
import { AuthForm } from "@/components/auth-form";
export const metadata = { title: "Вход" };

// Page отображает содержимое маршрута.
export default function Page() {
  return <AuthForm />;
}
