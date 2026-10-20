// Личный кабинет с защищённым чтением профиля и заявок.
import type { Metadata } from "next";
import { redirect } from "next/navigation";

import { Dashboard } from "@/components/dashboard";
import { ServiceUnavailable } from "@/components/service-unavailable";
import type { CleaningRequest, CleaningService } from "@/lib/api";
import { getSession, readAPI } from "@/lib/server-api";

export const metadata: Metadata = { title: "Личный кабинет" };

/** Параметры компонента DashboardPageProps. */
interface DashboardPageProps {
  searchParams: Promise<{ service?: string }>;
}

/** Загружает собственные заявки и обрабатывает недоступность API. */
export default async function DashboardPage({
  searchParams,
}: DashboardPageProps) {
  const parameters = await searchParams;
  const selectedService =
    typeof parameters.service === "string" ? parameters.service : "";
  const serviceQuery = selectedService
    ? `?service=${encodeURIComponent(selectedService)}`
    : "";
  const { user, unavailable } = await getSession();

  if (unavailable) {
    return <ServiceUnavailable />;
  }

  if (!user) {
    redirect(`/login${serviceQuery}`);
  }

  let catalog: { services: CleaningService[] };
  let requests: { requests: CleaningRequest[] };

  try {
    [catalog, requests] = await Promise.all([
      readAPI<{ services: CleaningService[] }>("/services"),
      readAPI<{ requests: CleaningRequest[] }>("/requests"),
    ]);
  } catch {
    return <ServiceUnavailable />;
  }

  return (
    <Dashboard
      user={user}
      services={catalog.services}
      initialRequests={requests.requests}
      selectedService={selectedService}
    />
  );
}
