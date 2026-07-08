import { useMutation } from "@tanstack/react-query";
import { login as apiLogin, register as apiRegister } from "../api/auth";
import {
  createMyProvider,
  type ProviderInput,
} from "../api/providers";
import { setUserInterests } from "../api/tags";
import { useAuthStore } from "../stores/auth.store";

export function useLogin() {
  const setAuth = useAuthStore((s) => s.setAuth);

  return useMutation({
    mutationFn: ({
      email,
      password,
    }: {
      email: string;
      password: string;
    }) => apiLogin(email, password),
    onSuccess: (data) => setAuth(data.user, data.access_token),
  });
}

export function useRegister() {
  const setAuth = useAuthStore((s) => s.setAuth);

  return useMutation({
    mutationFn: ({
      name,
      email,
      password,
      role = "customer",
    }: {
      name: string;
      email: string;
      password: string;
      role?: string;
    }) => apiRegister(name, email, password, role),
    onSuccess: (data) => setAuth(data.user, data.access_token),
  });
}

export function useCreateProviderProfile() {
  return useMutation({
    mutationFn: (data: ProviderInput) => createMyProvider(data),
  });
}

export function useSetUserInterests() {
  return useMutation({
    mutationFn: (tagIds: number[]) => setUserInterests(tagIds),
  });
}
