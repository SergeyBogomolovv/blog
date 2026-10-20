// Форма оформления заявки на уборку.
"use client";

import { useRouter } from "next/navigation";
import { type FormEvent, useState } from "react";
import { ArrowUpRight, Check, Plus } from "lucide-react";

import {
  api,
  ApiError,
  /** Структура данных CleaningRequest. */
  type CleaningRequest,
  /** Структура данных CleaningService. */
  type CleaningService,
  errorMessage,
  /** Структура данных FieldErrors. */
  type FieldErrors,
  moscowToday,
} from "@/lib/api";
import { validateRequest } from "@/lib/validation";

/** Параметры компонента RequestFormProps. */
interface RequestFormProps {
  services: CleaningService[];
  selectedService: string;
  onCreated: (request: CleaningRequest) => void;
}

/** Собирает услугу, адрес и дату уборки с проверкой полей. */
export function RequestForm({
  services,
  selectedService,
  onCreated,
}: RequestFormProps) {
  const router = useRouter();
  const initialService = services.some(
    (service) => service.id === selectedService,
  )
    ? selectedService
    : "";

  const [serviceID, setServiceID] = useState(initialService);
  const [address, setAddress] = useState("");
  const [cleaningDate, setCleaningDate] = useState("");
  const [pending, setPending] = useState(false);
  const [errors, setErrors] = useState<FieldErrors>({});
  const [error, setError] = useState("");
  const [success, setSuccess] = useState(false);

  /** Проверяет поля и отправляет заявку. */
  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setSuccess(false);

    const input = {
      service: serviceID,
      address: address.trim(),
      date: cleaningDate,
    };
    const validationErrors = validateRequest(input, services, moscowToday());
    setErrors(validationErrors);

    if (Object.keys(validationErrors).length > 0) {
      return;
    }

    setPending(true);

    try {
      const response = await api<{ request: CleaningRequest }>(
        "/requests",
        input,
      );
      onCreated(response.request);
      setAddress("");
      setCleaningDate("");
      setSuccess(true);
    } catch (error) {
      if (error instanceof ApiError) {
        setErrors(error.fields);

        if (error.status === 401) {
          router.push("/login");
          router.refresh();

          return;
        }
      }

      setError(errorMessage(error));
    } finally {
      setPending(false);
    }
  }

  return (
    <section
      className="booking-panel"
      aria-labelledby="booking-title"
    >
      <span className="panel-icon">
        <Plus size={22} />
      </span>
      <h2 id="booking-title">Запланировать уборку</h2>
      <p>Выберите услугу и расскажите, где нас ждут.</p>
      <form
        method="post"
        className="form-stack booking-form"
        onSubmit={submit}
        autoComplete="off"
      >
        {error && (
          <div
            className="error-box"
            role="alert"
          >
            {error}
          </div>
        )}
        {success && (
          <div
            className="success-box"
            role="status"
          >
            <Check size={18} /> Заявка сохранена и появилась в вашем списке.
          </div>
        )}
        <div className="field">
          <label htmlFor="service">Вид уборки</label>
          <select
            id="service"
            name="service"
            value={serviceID}
            onChange={(event) => setServiceID(event.target.value)}
            required
            disabled={pending}
            aria-invalid={Boolean(errors.service)}
            aria-describedby="service-error"
          >
            <option
              value=""
              disabled
            >
              Выберите подходящую услугу
            </option>
            {services.map((service) => (
              <option
                key={service.id}
                value={service.id}
              >
                {service.name}
              </option>
            ))}
          </select>
          <span
            className="field-error"
            id="service-error"
          >
            {errors.service}
          </span>
        </div>
        <div className="field">
          <label htmlFor="address">Адрес уборки</label>
          <textarea
            id="address"
            name="address"
            rows={3}
            required
            maxLength={100}
            value={address}
            placeholder="Город, улица, дом, квартира"
            disabled={pending}
            onChange={(event) => setAddress(event.target.value)}
            onBlur={() => setAddress((current) => current.trim())}
            aria-invalid={Boolean(errors.address)}
            aria-describedby="address-hint address-error"
          />
          <div
            className="field-hint field-counter"
            id="address-hint"
          >
            <span>Полный адрес, до 100 символов</span>
            <span>{Array.from(address).length}/100</span>
          </div>
          <span
            className="field-error"
            id="address-error"
          >
            {errors.address}
          </span>
        </div>
        <div className="field">
          <label htmlFor="cleaning-date">Дата уборки</label>
          <input
            id="cleaning-date"
            name="date"
            type="date"
            min={moscowToday()}
            required
            value={cleaningDate}
            onChange={(event) => setCleaningDate(event.target.value)}
            disabled={pending}
            aria-invalid={Boolean(errors.date)}
            aria-describedby="date-hint date-error"
          />
          <span
            className="field-hint"
            id="date-hint"
          >
            Сегодня или позже. Часовой пояс — Москва.
          </span>
          <span
            className="field-error"
            id="date-error"
          >
            {errors.date}
          </span>
        </div>
        <button
          type="submit"
          className="button wide-button"
          disabled={pending}
        >
          {pending ? "Сохраняем заявку…" : "Отправить заявку"}
          <ArrowUpRight size={18} />
        </button>
        <span className="field-hint">
          После отправки заявка сохранится в вашем кабинете.
        </span>
      </form>
    </section>
  );
}
