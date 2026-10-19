// Общие типы ответов Go API и форматирование дат без сетевых запросов.
export interface User {
  id: number;
  login: string;
  name: string;
  bio: string;
  email?: string;
  role: "user" | "admin";
  banned: boolean;
  created_at: string;
}

export interface Post {
  id: number;
  author_id: number;
  title: string;
  body: string;
  created_at: string;
  comment_count: number;
  author: User;
}

export interface Comment {
  id: number;
  post_id: number;
  body: string;
  created_at: string;
  author: User;
}

// date форматирует дату одинаково на сервере и в браузере.
export function date(value: string) {
  return new Intl.DateTimeFormat("ru-RU", {
    day: "numeric",
    month: "long",
    year: "numeric",
    timeZone: "Europe/Moscow",
  }).format(new Date(value));
}
