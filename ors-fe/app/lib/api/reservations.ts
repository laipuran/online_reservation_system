import { request } from "./client";

export type ReservationStatus =
  | "pending"
  | "confirmed"
  | "completed"
  | "cancelled"
  | "rejected";

export interface ReservationItem {
  id: number;
  user_id: number;
  service_id: number;
  start_time: string;
  end_time: string;
  status: ReservationStatus;
  note?: string;
  created_at: string;
  updated_at: string;
}

export interface ReservationViewItem {
  id: number;
  service: {
    id: number;
    title: string;
    provider: {
      id: number;
      business_name: string;
    };
  };
  start_time: string;
  end_time: string;
  status: ReservationStatus;
  note?: string;
  created_at: string;
}

export interface ReservationListResponse {
  items: ReservationItem[];
  page: number;
  page_size: number;
}

export interface ReservationQueryParams {
  status?: ReservationStatus;
  page?: number;
  page_size?: number;
}

/* ── Customer ────────────────────────────────────────── */

export function fetchMyReservations(
  params: ReservationQueryParams = {}
): Promise<ReservationListResponse> {
  const query = new URLSearchParams();
  if (params.status) query.set("status", params.status);
  if (params.page) query.set("page", String(params.page));
  if (params.page_size) query.set("page_size", String(params.page_size));
  const qs = query.toString();
  return request<ReservationListResponse>(
    `/reservations${qs ? `?${qs}` : ""}`
  );
}

export function fetchReservationById(
  id: number
): Promise<ReservationItem> {
  return request<ReservationItem>(`/reservations/${id}`);
}

export interface CreateReservationInput {
  service_id: number;
  start_time: string;
  note?: string;
}

export function createReservation(
  data: CreateReservationInput
): Promise<ReservationViewItem> {
  return request<ReservationViewItem>("/reservations", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export function cancelReservation(id: number): Promise<ReservationItem> {
  return request<ReservationItem>(`/reservations/${id}/cancel`, {
    method: "PUT",
  });
}

/* ── Provider ────────────────────────────────────────── */

export function fetchProviderReservations(
  params: ReservationQueryParams = {}
): Promise<ReservationListResponse> {
  const query = new URLSearchParams();
  if (params.status) query.set("status", params.status);
  if (params.page) query.set("page", String(params.page));
  if (params.page_size) query.set("page_size", String(params.page_size));

  const qs = query.toString();
  return request<ReservationListResponse>(
    `/provider/reservations${qs ? `?${qs}` : ""}`
  );
}

export function confirmReservation(id: number): Promise<ReservationItem> {
  return request<ReservationItem>(`/provider/reservations/${id}/confirm`, {
    method: "PUT",
  });
}

export function rejectReservation(id: number): Promise<ReservationItem> {
  return request<ReservationItem>(`/provider/reservations/${id}/reject`, {
    method: "PUT",
  });
}
