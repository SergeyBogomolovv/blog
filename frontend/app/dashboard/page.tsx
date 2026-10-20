import type { Metadata } from "next";
import { redirect } from "next/navigation";

import { Dashboard } from "@/components/dashboard";
import type { CleaningRequest, CleaningService } from "@/lib/api";
import { getSession, readAPI } from "@/lib/server-api";

export const metadata: Metadata = { title: "Личный кабинет" };

interface DashboardPageProps {
  searchParams: Promise<{ service?: string }>;
}

export default async function DashboardPage({ searchParams }: DashboardPageProps) {
  const parameters = await searchParams;
  const selectedService = typeof parameters.service === "string"
    ? parameters.service
    : "";
  const serviceQuery = selectedService
    ? `?service=${encodeURIComponent(selectedService)}`
    : "";
  const user = await getSession();

  if (!user) {
    redirect(`/login${serviceQuery}`);
  }

  const [catalog, requests] = await Promise.all([
    readAPI<{ services: CleaningService[] }>("/services"),
    readAPI<{ requests: CleaningRequest[] }>("/requests"),
  ]);

  return (
    <Dashboard
      user={user}
      services={catalog.services}
      initialRequests={requests.requests}
      selectedService={selectedService}
    />
  );
}
