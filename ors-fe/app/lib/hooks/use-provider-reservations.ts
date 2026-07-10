import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  fetchProviderReservations,
  confirmReservation,
  rejectReservation,
  type ReservationQueryParams,
  type ReservationViewItem,
} from "../api/reservations";
import { fetchServiceById } from "../api/services";

export function useProviderReservations(
  params: ReservationQueryParams = {}
) {
  return useQuery({
    queryKey: ["provider-reservations", params],
    queryFn: async () => {
      const result = await fetchProviderReservations(params);

      const serviceIds = [...new Set(result.items.map((item) => item.service_id))];
      const services = await Promise.all(
        serviceIds.map((id) => fetchServiceById(id).catch(() => null))
      );
      const serviceMap = new Map(services.filter(Boolean).map((s) => [s!.id, s!]));

      const items: ReservationViewItem[] = result.items.map((item) => {
        const svc = serviceMap.get(item.service_id);
        return {
          id: item.id,
          user_id: item.user_id,
          service: {
            id: svc?.id ?? item.service_id,
            title: svc?.title ?? "未知服务",
            provider: {
              id: svc?.provider?.id ?? 0,
              business_name: svc?.provider?.business_name ?? "未知商家",
            },
          },
          start_time: item.start_time,
          end_time: item.end_time,
          status: item.status,
          note: item.note,
          created_at: item.created_at,
        };
      });

      return { items, page: result.page, page_size: result.page_size };
    },
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
