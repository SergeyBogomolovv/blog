// Серверное отображение инициалов автора без внешних изображений.
// Avatar формирует одинаковый знак автора в ленте, профиле и обсуждениях.
export function Avatar({
  name,
  large = false,
}: {
  name: string;
  large?: boolean;
}) {
  return (
    <span
      className={`avatar${large ? " large" : ""}`}
      aria-hidden="true"
    >
      {name.trim().slice(0, 1).toUpperCase() || "?"}
    </span>
  );
}
