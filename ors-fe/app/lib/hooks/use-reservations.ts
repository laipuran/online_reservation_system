import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  fetchMyReservations,
  cancelReservation,
  type ReservationQueryParams,
  type ReservationViewItem,
} from "../api/reservations";
import { fetchServiceById } from "../api/services";

export function useMyReservations(
  params: ReservationQueryParams = {}
) {
  return useQuery({
    queryKey: ["my-reservations", params],
    queryFn: async () => {
      const result = await fetchMyReservations(params);

      const serviceIds = [...new Set(result.items.map((item) => item.service_id))];
      const services = await Promise.all(
        serviceIds.map((id) => fetchServiceById(id).catch(() => null))
      );
      const serviceMap = new Map(services.filter(Boolean).map((s) => [s!.id, s!]));

      const items: ReservationViewItem[] = result.items.map((item) => {
        const svc = serviceMap.get(item.service_id);
        return {
          id: item.id,
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

export function useCancelReservation() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => cancelReservation(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["my-reservations"] });
    },
  });
}
