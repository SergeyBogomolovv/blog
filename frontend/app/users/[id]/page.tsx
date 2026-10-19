// Серверная страница пользователя с пагинацией его публикаций.
import { Profile } from "@/components/profile";
import { pageNumber, resourceID } from "@/lib/params";

// Page проверяет параметры до запроса профиля и постов.
export default async function Page({
  params,
  searchParams,
}: {
  params: Promise<{ id: string }>;
  searchParams: Promise<{ page?: string | string[] }>;
}) {
  const [route, query] = await Promise.all([params, searchParams]);

  return (
    <Profile
      id={resourceID(route.id)}
      page={pageNumber(query.page)}
    />
  );
}
