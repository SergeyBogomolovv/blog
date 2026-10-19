// Страница: Регистрация.
import { AuthForm } from "@/components/auth-form";
export const metadata = { title: "Регистрация" };

// Page отображает содержимое маршрута.
export default function Page() {
  return <AuthForm register />;
}
