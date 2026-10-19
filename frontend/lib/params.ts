// Проверка недоверенных параметров маршрутов перед обращением к API.
import { notFound } from "next/navigation";

// pageNumber нормализует номер страницы и ограничивает доступный диапазон.
export function pageNumber(value: string | string[] | undefined): number {
  if (value === undefined) {
    return 0;
  }

  if (
    typeof value !== "string" ||
    !/^\d{1,5}$/.test(value) ||
    Number(value) > 10000
  ) {
    notFound();
  }

  return Number(value);
}

// resourceID разрешает только положительные числовые идентификаторы маршрута.
export function resourceID(value: string): string {
  if (!/^[1-9]\d{0,15}$/.test(value)) {
    notFound();
  }

  return value;
}
