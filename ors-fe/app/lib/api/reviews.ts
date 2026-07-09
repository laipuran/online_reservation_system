import { request } from "./client";

export interface ReviewStats {
  avg_rating: number;
  total: number;
  distribution: Record<number, number>;
}

export interface ReviewItem {
  id: number;
  reservation_id: number;
  user_id: number;
  user_name: string;
  service_id: number;
  rating: number;
  comment: string;
  created_at: string;
}

export interface ReviewListResponse {
  items: ReviewItem[];
  total: number;
  page: number;
  page_size: number;
}

export function fetchServiceReviewStats(
  serviceId: number
): Promise<ReviewStats> {
  return request<ReviewStats>(`/services/${serviceId}/reviews/stats`);
}

export function fetchServiceReviews(
  serviceId: number,
  params?: { page?: number; page_size?: number }
): Promise<ReviewListResponse> {
  const query = new URLSearchParams();
  if (params?.page) query.set("page", String(params.page));
  if (params?.page_size) query.set("page_size", String(params.page_size));
  const qs = query.toString();
  return request<ReviewListResponse>(
    `/services/${serviceId}/reviews${qs ? `?${qs}` : ""}`
  );
}
