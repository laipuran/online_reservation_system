import { useQuery } from "@tanstack/react-query";
import { fetchServiceReviews } from "../api/reviews";

export function useServiceReviews(
  serviceId: number,
  params?: { page?: number; page_size?: number }
) {
  return useQuery({
    queryKey: ["service-reviews", serviceId, params],
    queryFn: () => fetchServiceReviews(serviceId, params),
    enabled: !!serviceId,
  });
}
