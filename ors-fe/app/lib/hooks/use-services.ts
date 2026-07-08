import { useQuery } from "@tanstack/react-query";
import { fetchServices, type ServiceQueryParams } from "../api/services";

export function useServices(params: ServiceQueryParams = {}) {
  return useQuery({
    queryKey: ["services", params],
    queryFn: () => fetchServices(params),
  });
}
