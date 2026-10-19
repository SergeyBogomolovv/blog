// Состояние ожидания потокового серверного рендера маршрута.
// Loading показывается Next.js при переходе до готовности данных.
export default function Loading() {
  return (
    <p
      className="loading"
      role="status"
    >
      Открываем страницу…
    </p>
  );
}
