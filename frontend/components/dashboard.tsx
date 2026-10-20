"use client";

import Link from "next/link";
import { useState } from "react";
import { ShieldCheck, Sparkles } from "lucide-react";

import { RequestForm } from "@/components/request-form";
import { RequestsList } from "@/components/requests-list";
import {
  type CleaningRequest,
  type CleaningService,
  date,
  type User,
} from "@/lib/api";

interface DashboardProps {
  user: User;
  services: CleaningService[];
  initialRequests: CleaningRequest[];
  selectedService: string;
}

function AccountDetails({ user }: { user: User }) {
  return (
    <div className="account-details">
      <Sparkles size={19} />
      <div>
        <strong>{user.surname} {user.name} {user.patronymic}</strong>
        <p>{user.phone} · {user.email}</p>
        <Link href="/privacy">
          Согласие от {date(user.consent_at)}, версия {user.consent_version}
        </Link>
      </div>
    </div>
  );
}

export function Dashboard({
  user,
  services,
  initialRequests,
  selectedService,
}: DashboardProps) {
  const [requests, setRequests] = useState(initialRequests);

  function addRequest(request: CleaningRequest) {
    setRequests((current) => [request, ...current]);
  }

  return (
    <div className="container dashboard">
      <div className="dashboard-heading">
        <div>
          <span className="eyebrow">Ваше пространство заботы</span>
          <h1>Здравствуйте, <em>{user.name}.</em></h1>
          <p>Здесь начинается порядок — и хранятся все ваши заявки.</p>
        </div>
        <span className="account-badge">
          <ShieldCheck size={17} /> Личный кабинет
        </span>
      </div>
      <div className="dashboard-grid">
        <RequestForm
          services={services}
          selectedService={selectedService}
          onCreated={addRequest}
        />
        <section className="requests-panel" aria-labelledby="requests-title">
          <div className="requests-heading">
            <h2 id="requests-title">Мои заявки <span>{requests.length}</span></h2>
            <span>Сначала новые</span>
          </div>
          <RequestsList requests={requests} services={services} />
          <AccountDetails user={user} />
        </section>
      </div>
    </div>
  );
}
