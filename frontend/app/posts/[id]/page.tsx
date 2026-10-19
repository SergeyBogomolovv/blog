// Серверная страница публикации с настоящим 404 для отсутствующего поста.
import { PostDetail } from "@/components/post-detail";
import { resourceID } from "@/lib/params";

// Page проверяет идентификатор маршрута перед запросом к Go API.
export default async function Page({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  return <PostDetail id={resourceID((await params).id)} />;
}
