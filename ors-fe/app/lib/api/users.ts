import { request, type User } from "./client";

export function fetchMyProfile(): Promise<User> {
  return request<User>("/users/me");
}

export interface UpdateProfileInput {
  name?: string;
  phone?: string;
  avatar_url?: string;
}

export function updateMyProfile(data: UpdateProfileInput): Promise<User> {
  return request<User>("/users/me", {
    method: "PUT",
    body: JSON.stringify(data),
  });
}

export interface UpdatePasswordInput {
  current_password: string;
  new_password: string;
}

export function updateMyPassword(data: UpdatePasswordInput): Promise<{ message: string }> {
  return request<{ message: string }>("/users/me/password", {
    method: "PUT",
    body: JSON.stringify(data),
  });
}

export function fetchMyInterests(): Promise<number[]> {
  return request<number[]>("/users/me/interests");
}
