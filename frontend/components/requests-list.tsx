import {
  ArrowUpRight,
  CalendarDays,
  Check,
  ClipboardList,
  MapPin,
} from "lucide-react";

import { type CleaningRequest, type CleaningService, date } from "@/lib/api";

interface RequestsListProps {
  requests: CleaningRequest[];
  services: CleaningService[];
}

function EmptyRequests() {
  return (
    <div className="empty-requests">
      <span><ClipboardList size={34} strokeWidth={1.3} /></span>
      <h3>Здесь пока чистый лист</h3>
      <p>
        Запланируйте первую уборку.
        <br />
        Ваша заявка появится здесь сразу после отправки.
      </p>
      <a href="#service" className="text-link">
        Начать с выбора услуги <ArrowUpRight size={16} />
      </a>
    </div>
  );
}

function RequestCard({
  request,
  serviceName,
}: {
  request: CleaningRequest;
  serviceName: string;
}) {
  return (
    <article className="request-card">
      <div className="request-top">
        <span>Заявка № {request.id}</span>
        <span className="saved-label"><Check size={13} /> Сохранена</span>
      </div>
      <h3>{serviceName}</h3>
      <p><MapPin size={16} />{request.address}</p>
      <p><CalendarDays size={16} />{date(request.cleaning_date)}</p>
      <div className="request-bottom">Создана {date(request.created_at)}</div>
    </article>
  );
}

export function RequestsList({ requests, services }: RequestsListProps) {
  if (requests.length === 0) {
    return <EmptyRequests />;
  }

  const serviceNames = new Map(
    services.map((service) => [service.id, service.name]),
  );

  return (
    <div className="requests-list" aria-live="polite">
      {requests.map((request) => (
        <RequestCard
          key={request.id}
          request={request}
          serviceName={serviceNames.get(request.service) || request.service}
        />
      ))}
    </div>
  );
}
