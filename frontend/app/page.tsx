// Главная страница с каталогом услуг и описанием оформления заявки.
import Link from "next/link";
import {
  ArrowDown,
  ArrowUpRight,
  Check,
  Leaf,
  ShieldCheck,
  Sparkles,
} from "lucide-react";
import { Services } from "@/components/services";
import { readAPI } from "@/lib/server-api";
import type { CleaningService } from "@/lib/api";

const bookingSteps = [
  {
    title: "Создайте аккаунт",
    text: "Несколько контактных данных — и у вас появится свой личный кабинет.",
  },
  {
    title: "Расскажите о задаче",
    text: "Выберите вид уборки, укажите адрес и удобную дату.",
  },
  {
    title: "Сохраните заявку",
    text: "Она сразу появится в кабинете. Все ваши обращения будут под рукой.",
  },
];

/** Загружает каталог и отображает главную страницу сервиса. */
export default async function Home() {
  let services: CleaningService[] = [];
  let error = "";

  try {
    const response = await readAPI<{ services: CleaningService[] }>(
      "/services",
    );
    services = response.services;
  } catch {
    error =
      "Не удалось загрузить услуги. Проверьте запуск сервера и обновите страницу.";
  }

  return (
    <>
      <link
        rel="preload"
        as="image"
        href="/clean-room.svg"
      />
      <section className="hero container">
        <div className="hero-copy">
          <span className="eyebrow">
            <span className="status-dot" /> Чистота, с которой легче жить
          </span>
          <h1>
            Пусть дома
            <br />
            будет <em>хорошо.</em>
          </h1>
          <p className="hero-description">
            Возьмём уборку на себя, чтобы у вас оставалось больше времени на то,
            что вы любите.
          </p>
          <Link
            href="/dashboard"
            className="button"
          >
            Запланировать уборку <ArrowUpRight size={19} />
          </Link>
          <a
            href="#services"
            className="hero-secondary"
          >
            Подобрать услугу <ArrowDown size={15} />
          </a>
          <div className="hero-note">
            <span className="note-icon">
              <Leaf size={20} />
            </span>
            <span>
              Для больших планов
              <br />
              <strong>и маленьких домашних радостей</strong>
            </span>
          </div>
        </div>
        <div className="hero-art">
          <div className="art-label">
            <Sparkles size={17} /> Место, где хочется быть
          </div>
          <div
            className="room-image"
            role="img"
            aria-label="Светлая гостиная с зелёным диваном, растениями и большим солнечным окном"
          />
          <div className="art-stamp">
            <ShieldCheck size={27} />
            <span>
              Ваш дом.
              <br />
              <strong>В заботливых руках.</strong>
            </span>
            <Check size={17} />
          </div>
          <span className="art-caption">ЧИСТЫЙ ДОМ — СПОКОЙНАЯ ГОЛОВА</span>
        </div>
      </section>
      <div className="values-strip">
        <div className="container values-inner">
          <span>
            <Check size={17} /> Три вида уборки
          </span>
          <span>
            <Check size={17} /> Удобная дата
          </span>
          <span>
            <Check size={17} /> Заявка за пару минут
          </span>
          <span>
            <Check size={17} /> Всё в личном кабинете
          </span>
        </div>
      </div>
      <section
        className="container section"
        id="services"
      >
        <div className="section-heading">
          <div>
            <span className="eyebrow">01 / Наши услуги</span>
            <h2>
              Каждому дому —<br />
              <em>своя забота.</em>
            </h2>
          </div>
          <p>
            Освежить пространство, навести основательный порядок
            <br className="desktop-break" /> или начать с чистого листа после
            ремонта.
          </p>
        </div>
        <Services
          initialServices={services}
          initialError={error}
        />
      </section>
      <section
        className="how-section"
        id="how-it-works"
      >
        <div className="container section">
          <div className="section-heading">
            <div>
              <span className="eyebrow">02 / Всё просто</span>
              <h2>
                Чистота начинается
                <br />
                <em>с трёх шагов.</em>
              </h2>
            </div>
            <Link
              className="text-link"
              href="/register"
            >
              Начать сейчас <ArrowUpRight size={17} />
            </Link>
          </div>
          <div className="steps">
            {bookingSteps.map((step, index) => (
              <div
                className="step"
                key={step.title}
              >
                <span className="step-number">0{index + 1}</span>
                <h3>{step.title}</h3>
                <p>{step.text}</p>
              </div>
            ))}
          </div>
        </div>
      </section>
      <section className="container bottom-cta">
        <Sparkles size={36} />
        <div>
          <h2>
            Освободите время <em>для себя.</em>
          </h2>
          <p>А первый шаг к чистому дому сделайте прямо сейчас.</p>
        </div>
        <Link
          href="/dashboard"
          className="button"
        >
          Оставить заявку <ArrowUpRight size={18} />
        </Link>
      </section>
    </>
  );
}
