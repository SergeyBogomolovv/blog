// Каталог услуг и поиск с проверкой пользовательского запроса.
"use client";

import Link from "next/link";
import { type FormEvent, useState } from "react";
import {
  ArrowUpRight,
  Brush,
  /** Структура данных LucideIcon. */
  type LucideIcon,
  Search,
  Sparkles,
  SprayCan,
  X,
} from "lucide-react";

import { api, type CleaningService, errorMessage } from "@/lib/api";
import { validSearch } from "@/lib/validation";

/** Структура данных ServiceAppearance. */
interface ServiceAppearance {
  icon: LucideIcon;
  caption: string;
  number: string;
}

const serviceAppearances: Record<string, ServiceAppearance> = {
  maintenance: {
    icon: SprayCan,
    caption: "Для повседневного комфорта",
    number: "01",
  },
  deep: {
    icon: Sparkles,
    caption: "Внимание к каждой детали",
    number: "02",
  },
  renovation: {
    icon: Brush,
    caption: "Когда всё самое сложное позади",
    number: "03",
  },
};

const defaultAppearance: ServiceAppearance = {
  icon: Sparkles,
  caption: "Забота о вашем доме",
  number: "01",
};

/** Отображает услугу и ссылку на оформление заявки. */
function ServiceCard({ service }: { service: CleaningService }) {
  const appearance = serviceAppearances[service.id] || defaultAppearance;
  const Icon = appearance.icon;

  return (
    <article className={`service-card service-${service.id}`}>
      <div className="card-top">
        <span className="service-icon">
          <Icon
            size={30}
            strokeWidth={1.4}
          />
        </span>
        <span className="card-number">/ {appearance.number}</span>
      </div>
      <span className="service-tag">{appearance.caption}</span>
      <h3>{service.name}</h3>
      <p>{service.description}</p>
      <Link
        href={`/dashboard?service=${service.id}`}
        className="card-link"
      >
        Выбрать уборку
        <span>
          <ArrowUpRight size={20} />
        </span>
      </Link>
    </article>
  );
}

/** Параметры компонента ServicesProps. */
interface ServicesProps {
  initialServices: CleaningService[];
  initialError: string;
}

/** Управляет каталогом, поиском и сообщениями об отсутствии результатов. */
export function Services({ initialServices, initialError }: ServicesProps) {
  const [services, setServices] = useState(initialServices);
  const [query, setQuery] = useState("");
  const [appliedQuery, setAppliedQuery] = useState("");
  const [error, setError] = useState(initialError);
  const [pending, setPending] = useState(false);

  /** Проверяет запрос и обновляет список найденных услуг. */
  async function loadServices(searchQuery: string) {
    const normalizedQuery = searchQuery.trim();

    if (normalizedQuery && !validSearch(normalizedQuery)) {
      setError(
        "До 20 символов: русские и латинские буквы, цифры, пробелы и дефисы.",
      );

      return;
    }

    setPending(true);
    setError("");

    try {
      const parameters = new URLSearchParams();

      if (normalizedQuery) {
        parameters.set("q", normalizedQuery);
      }

      const suffix = normalizedQuery ? `?${parameters}` : "";
      const response = await api<{ services: CleaningService[] }>(
        `/services${suffix}`,
      );
      setServices(response.services);
      setAppliedQuery(normalizedQuery);
    } catch (error) {
      setError(errorMessage(error));
    } finally {
      setPending(false);
    }
  }

  /** Обрабатывает отправку поисковой формы без перезагрузки страницы. */
  function submitSearch(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    void loadServices(query);
  }

  /** Сбрасывает запрос и загружает полный каталог услуг. */
  function resetSearch() {
    setQuery("");
    void loadServices("");
  }

  const noResults = !error && services.length === 0;

  return (
    <>
      <form
        className="search-form"
        onSubmit={submitSearch}
        role="search"
        autoComplete="off"
      >
        <Search size={20} />
        <label
          className="sr-only"
          htmlFor="service-search"
        >
          Поиск услуг
        </label>
        <input
          id="service-search"
          value={query}
          onChange={(event) => setQuery(event.target.value)}
          maxLength={20}
          placeholder="Какая уборка вам нужна?"
          aria-describedby={error ? "search-error" : "search-hint"}
        />
        {query && (
          <button
            type="button"
            className="icon-button"
            aria-label="Сбросить поиск"
            onClick={resetSearch}
            disabled={pending}
          >
            <X size={18} />
          </button>
        )}
        <button
          type="submit"
          className="search-submit"
          disabled={pending}
        >
          {pending ? "Ищем…" : "Найти услугу"}
          <ArrowUpRight size={16} />
        </button>
      </form>
      <p
        className="search-hint"
        id="search-hint"
      >
        Например, «ремонта». Буквы, цифры, пробелы и дефисы — до 20 символов.
      </p>
      {error && (
        <div
          className="error-box"
          role="alert"
          id="search-error"
        >
          {error}{" "}
          <button
            className="text-link"
            onClick={resetSearch}
            disabled={pending}
          >
            Загрузить все услуги
          </button>
        </div>
      )}
      {noResults ? (
        <div className="empty-state">
          <Search size={30} />
          <h3>Ничего не нашлось</h3>
          <p>По запросу «{appliedQuery}» услуг нет. Попробуйте другое слово.</p>
          <button
            className="button button-outline"
            onClick={resetSearch}
          >
            Показать все услуги
          </button>
        </div>
      ) : (
        <div
          className="service-grid"
          aria-live="polite"
        >
          {services.map((service) => (
            <ServiceCard
              key={service.id}
              service={service}
            />
          ))}
        </div>
      )}
    </>
  );
}
