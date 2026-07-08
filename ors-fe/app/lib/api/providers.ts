import { request } from "./client";

export interface ProviderDetail {
  id: number;
  user_id: number;
  business_name: string;
  description: string;
  address: string;
  phone: string;
  email: string;
  logo_url: string;
  created_at: string;
  updated_at: string;
}

export interface ProviderInput {
  business_name: string;
  description?: string;
  address?: string;
  phone?: string;
  email?: string;
  logo_url?: string;
}

export function fetchProvider(id: number): Promise<ProviderDetail> {
  return request<ProviderDetail>(`/providers/${id}`);
}

export function createMyProvider(data: ProviderInput): Promise<ProviderDetail> {
  return request<ProviderDetail>("/providers/me", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export function fetchMyProvider(): Promise<ProviderDetail> {
  return request<ProviderDetail>("/providers/me");
}

export function updateMyProvider(data: ProviderInput): Promise<ProviderDetail> {
  return request<ProviderDetail>("/providers/me", {
    method: "PUT",
    body: JSON.stringify(data),
  });
}
