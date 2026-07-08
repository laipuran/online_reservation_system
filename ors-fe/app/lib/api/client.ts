import { useAuthStore } from "../stores/auth.store";

const BASE_URL = "http://localhost:8080/api/v1";

export interface ApiResponse<T = unknown> {
  code: number;
  message: string;
  data: T;
}

export interface User {
  id: number;
  name: string;
  email: string;
  role: string;
  phone: string;
  avatar_url: string;
  created_at: string;
  updated_at: string;
}

export interface AuthData {
  user: User;
  access_token: string;
}

export class ApiError extends Error {
  constructor(
    public status: number,
    public code: number,
    message: string
  ) {
    super(message);
  }
}

export async function request<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const token = useAuthStore.getState().token;
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...((options.headers as Record<string, string>) || {}),
  };
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(`${BASE_URL}${path}`, { ...options, headers });
  const body: ApiResponse<T> = await res.json();

  if (!res.ok) {
    throw new ApiError(res.status, body.code, body.message);
  }

  return body.data;
}
