import { useEffect, useState } from "react";
import { useAuthStore } from "../stores/auth.store";

export function useAuth() {
  const [hydrated, setHydrated] = useState(false);
  const user = useAuthStore((s) => s.user);
  const token = useAuthStore((s) => s.token);
  const setAuth = useAuthStore((s) => s.setAuth);
  const clearAuth = useAuthStore((s) => s.clearAuth);

  useEffect(() => {
    const unsub = useAuthStore.persist.onFinishHydration(() =>
      setHydrated(true)
    );
    if (useAuthStore.persist.hasHydrated()) {
      setHydrated(true);
    }
    return () => unsub();
  }, []);

  return { user, token, loading: !hydrated, setAuth, clearAuth };
}
