import { request } from "./client";

export type ServiceStatus = "active" | "inactive" | "pending" | "rejected";

export interface ServiceProvider {
  id: number;
  business_name: string;
}

export interface ServiceCategory {
  id: number;
  name: string;
}

export interface ServiceItem {
  id: number;
  title: string;
  description: string;
  price: number;
  duration_minutes: number;
  avg_rating: number;
  review_count: number;
  status: ServiceStatus;
  image_url: string;
  provider: ServiceProvider;
  category: ServiceCategory;
  created_at: string;
  updated_at: string;
}

export interface ServiceListResponse {
  items: ServiceItem[];
  total: number;
  page: number;
  page_size: number;
}

export interface ServiceQueryParams {
  keyword?: string;
  category_id?: number;
  provider_id?: number;
  min_price?: number;
  max_price?: number;
  sort_by?: "price" | "rating" | "created_at";
  sort_order?: "asc" | "desc";
  page?: number;
  page_size?: number;
}

export function fetchServices(
  params: ServiceQueryParams = {}
): Promise<ServiceListResponse> {
  const query = new URLSearchParams();
  if (params.keyword) query.set("keyword", params.keyword);
  if (params.category_id) query.set("category_id", String(params.category_id));
  if (params.provider_id) query.set("provider_id", String(params.provider_id));
  if (params.min_price) query.set("min_price", String(params.min_price));
  if (params.max_price) query.set("max_price", String(params.max_price));
  if (params.sort_by) query.set("sort_by", params.sort_by);
  if (params.sort_order) query.set("sort_order", params.sort_order);
  if (params.page) query.set("page", String(params.page));
  if (params.page_size) query.set("page_size", String(params.page_size));

  const qs = query.toString();
  return request<ServiceListResponse>(`/services${qs ? `?${qs}` : ""}`);
}

export function fetchServiceById(id: number): Promise<ServiceItem> {
  return request<ServiceItem>(`/services/${id}`);
}
