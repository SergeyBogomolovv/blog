// Защищённая серверная страница создания публикации.
import { redirect } from "next/navigation";

import { Editor } from "@/components/editor";
import { getSession } from "@/lib/server-api";
export const metadata = { title: "Новая история" };

// Page проверяет сессию до отображения формы редактора.
export default async function Page() {
  if (!(await getSession())) {
    redirect("/login");
  }

  return <Editor />;
}
