import { useMutation } from "@tanstack/react-query";
import { login as apiLogin, register as apiRegister } from "../api/auth";
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
    }: {
      name: string;
      email: string;
      password: string;
    }) => apiRegister(name, email, password),
    onSuccess: (data) => setAuth(data.user, data.access_token),
  });
}
