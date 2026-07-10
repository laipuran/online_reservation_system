import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  fetchServiceReviews,
  createReview as createReviewApi,
  type CreateReviewInput,
} from "../api/reviews";

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

export function useCreateReview() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateReviewInput) => createReviewApi(data),
    onSuccess: (_data, variables) => {
      qc.invalidateQueries({ queryKey: ["my-reservations"] });
      qc.invalidateQueries({ queryKey: ["service-reviews"] });
      qc.invalidateQueries({ queryKey: ["services"] });
      qc.invalidateQueries({ queryKey: ["service", variables.service_id] });
    },
  });
}
