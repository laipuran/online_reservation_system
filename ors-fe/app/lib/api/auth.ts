import { request, type AuthData } from "./client";

export function register(
  name: string,
  email: string,
  password: string,
  role: string = "customer"
): Promise<AuthData> {
  return request<AuthData>("/auth/register", {
    method: "POST",
    body: JSON.stringify({ name, email, password, role }),
  });
}
export function login(email: string, password: string): Promise<AuthData> {
  return request<AuthData>("/auth/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
}
