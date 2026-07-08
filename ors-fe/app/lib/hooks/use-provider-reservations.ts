import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  fetchProviderReservations,
  confirmReservation,
  rejectReservation,
  type ReservationQueryParams,
} from "../api/reservations";

export function useProviderReservations(
  params: ReservationQueryParams = {}
) {
  return useQuery({
    queryKey: ["provider-reservations", params],
    queryFn: () => fetchProviderReservations(params),
  });
}

export function useConfirmReservation() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => confirmReservation(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["provider-reservations"] });
    },
  });
}

export function useRejectReservation() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => rejectReservation(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["provider-reservations"] });
    },
  });
}
