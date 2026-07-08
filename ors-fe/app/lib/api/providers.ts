import { request } from "./client";

export interface ProviderProfileInput {
  business_name: string;
  description: string;
  address: string;
  email: string;
  phone?: string;
  logo_url?: string;
}

export interface ProviderProfile {
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

export function createProviderProfile(
  data: ProviderProfileInput
): Promise<ProviderProfile> {
  return request<ProviderProfile>("/providers/me", {
    method: "POST",
    body: JSON.stringify(data),
  });
}

export async function fetchMyProviderProfile(): Promise<ProviderProfile | null> {
  try {
    return await request<ProviderProfile>("/providers/me");
  } catch {
    return null;
  }
}

export function updateMyProviderProfile(
  data: ProviderProfileInput
): Promise<ProviderProfile> {
  return request<ProviderProfile>("/providers/me", {
    method: "PUT",
    body: JSON.stringify(data),
  });
}
