import { request } from "./client";

export interface CreateReviewInput {
  reservation_id: number;
  rating: number;
  comment?: string;
}

export interface ReviewItem {
  id: number;
  reservation_id: number;
  user_id: number;
  service_id: number;
  rating: number;
  comment?: string;
  created_at: string;
}

export interface ReviewListResponse {
  items: ReviewItem[];
  page: number;
  page_size: number;
}

export function createReview(data: CreateReviewInput): Promise<ReviewItem> {
  return request<ReviewItem>("/reviews", {
    method: "POST",
    body: JSON.stringify(data),
  });
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
