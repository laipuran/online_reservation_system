import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  fetchMyReservations,
  cancelReservation,
  type ReservationQueryParams,
} from "../api/reservations";

export function useMyReservations(
  params: ReservationQueryParams = {}
) {
  return useQuery({
    queryKey: ["my-reservations", params],
    queryFn: () => fetchMyReservations(params),
  });
}

export function useCancelReservation() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => cancelReservation(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["my-reservations"] });
    },
  });
}
