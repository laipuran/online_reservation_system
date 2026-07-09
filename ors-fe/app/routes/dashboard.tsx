import { useState, useEffect } from "react";
import { useNavigate, Link } from "react-router";
import { useQuery } from "@tanstack/react-query";
import { useAuth } from "../lib/hooks/use-auth";
import { fetchMyProvider } from "../lib/api/providers";
import { fetchServices, type ServiceItem } from "../lib/api/services";
import { useMyReservations, useCancelReservation } from "../lib/hooks/use-reservations";
import { useCreateReview } from "../lib/hooks/use-reviews";
import { ReservationCard } from "../lib/components/reservation-card";
import type { ReservationStatus } from "../lib/api/reservations";

const STATUS_FILTERS: Array<{ label: string; value: ReservationStatus | "" }> = [
  { label: "全部", value: "" },
  { label: "待确认", value: "pending" },
  { label: "已确认", value: "confirmed" },
  { label: "已完成", value: "completed" },
  { label: "已取消", value: "cancelled" },
  { label: "已拒绝", value: "rejected" },
];

const ROLE_LABEL: Record<string, string> = {
  customer: "服务体验者",
  provider: "服务提供者",
};

export default function Dashboard() {
  const { user, loading } = useAuth();
  const navigate = useNavigate();

  const [statusFilter, setStatusFilter] = useState<ReservationStatus | "">("");
  const [page, setPage] = useState(1);
  const pageSize = 10;

  const params = {
    ...(statusFilter ? { status: statusFilter as ReservationStatus } : {}),
    page,
    page_size: pageSize,
  };

  const { data: reservationsData, isLoading: reservationsLoading } = useMyReservations(params);
  const cancelMutation = useCancelReservation();
  const reviewMutation = useCreateReview();
  const [reviewedIds, setReviewedIds] = useState<number[]>([]);

  const { data: allServicesData } = useQuery({
    queryKey: ["services", { page_size: 50 }],
    queryFn: () => fetchServices({ page_size: 50 }),
  });

  const serviceMap = new Map<number, ServiceItem>();
  for (const s of allServicesData?.items ?? []) {
    serviceMap.set(s.id, s);
  }

  const reservations = (reservationsData?.items ?? []).map((r) => {
    const svc = serviceMap.get(r.service_id);
    return {
      ...r,
      service: svc ? { id: svc.id, title: svc.title, provider: { id: svc.provider.id, business_name: svc.provider.business_name } } : null,
    };
  });
  const hasMore = (reservationsData?.items?.length ?? 0) >= pageSize;

  const providerQuery = useQuery({
    queryKey: ["my-provider-profile"],
    queryFn: async () => {
      try {
        return await fetchMyProvider();
      } catch {
        return null;
      }
    },
    enabled: !loading && user?.role === "provider",
    retry: false,
  });

  useEffect(() => {
    if (!loading && !user) {
      navigate("/login", { replace: true });
    }
  }, [user, loading, navigate]);

  useEffect(() => {
    if (
      !loading &&
      user?.role === "provider" &&
      !providerQuery.isLoading &&
      providerQuery.data === null
    ) {
      navigate("/complete-profile", { replace: true });
    }
  }, [user, loading, providerQuery.isLoading, providerQuery.data, navigate]);

  if (loading || providerQuery.isLoading) {
    return (
      <div className="flex items-center justify-center mt-20">
        <p className="text-gray-500 dark:text-gray-400">加载中...</p>
      </div>
    );
  }

  if (!user) {
    return null;
  }

  const joinedDate = new Date(user.created_at).toLocaleDateString("zh-CN", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });

  return (
    <div className="max-w-3xl mx-auto mt-8 px-4 pb-12">
      {/* CustomerCard */}
      <div className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-5 flex items-center gap-4">
        <div className="w-12 h-12 rounded-full bg-blue-500 flex items-center justify-center text-white text-lg font-bold shrink-0">
          {user.name.charAt(0)}
        </div>
        <div className="min-w-0 flex-1">
          <h1 className="text-lg font-bold text-gray-900 dark:text-gray-100">{user.name}</h1>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
            {ROLE_LABEL[user.role] ?? user.role} · 注册于 {joinedDate}
          </p>
        </div>
        {user.role === "provider" && providerQuery.data && (
          <Link
            to="/provider/services"
            className="shrink-0 text-sm text-blue-600 dark:text-blue-400 border border-blue-200 dark:border-blue-800 rounded px-3 py-1.5 hover:bg-blue-50 dark:hover:bg-blue-900/20 transition-colors"
          >
            服务商控制台
          </Link>
        )}
      </div>

      {/* Status filter */}
      <div className="mt-6 flex items-center gap-2 overflow-x-auto pb-1">
        {STATUS_FILTERS.map((f) => (
          <button
            key={f.value}
            onClick={() => {
              setStatusFilter(f.value);
              setPage(1);
            }}
            className={`shrink-0 text-sm px-3 py-1.5 rounded-full border transition-colors ${
              statusFilter === f.value
                ? "bg-blue-600 text-white border-blue-600"
                : "bg-white dark:bg-gray-800 text-gray-600 dark:text-gray-400 border-gray-200 dark:border-gray-700 hover:border-blue-300"
            }`}
          >
            {f.label}
          </button>
        ))}
      </div>

      {/* Reservation cards */}
      <div className="mt-4 space-y-3">
        {reservationsLoading ? (
          <div className="flex justify-center py-16">
            <p className="text-gray-400 dark:text-gray-500">加载中...</p>
          </div>
        ) : reservations.length === 0 ? (
          <div className="flex justify-center py-16">
            <p className="text-gray-400 dark:text-gray-500">暂无预约记录</p>
          </div>
        ) : (
          reservations.map((r) => (
            <ReservationCard
              key={r.id}
              id={r.id}
              serviceTitle={r.service?.title ?? `服务 #${r.service_id}`}
              serviceId={r.service?.id ?? r.service_id}
              providerName={r.service?.provider.business_name ?? "未知商家"}
              startTime={r.start_time}
              endTime={r.end_time}
              status={r.status}
              note={r.note}
              onCancel={(id) => cancelMutation.mutate(id)}
              cancelPending={cancelMutation.isPending}
              onReview={(id, rating, comment) =>
                reviewMutation.mutate(
                  { reservation_id: id, rating, comment },
                  {
                    onSuccess: () => setReviewedIds((prev) => [...prev, id]),
                  }
                )
              }
              isReviewing={reviewMutation.isPending}
              reviewed={reviewedIds.includes(r.id)}
            />
          ))
        )}
      </div>

      {/* Pagination */}
      {reservations.length > 0 && (
        <div className="flex items-center justify-center gap-2 mt-6">
          <button
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            disabled={page <= 1}
            className="px-3 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded disabled:opacity-40 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
          >
            上一页
          </button>
          <span className="text-sm text-gray-500 dark:text-gray-400">第 {page} 页</span>
          <button
            onClick={() => setPage((p) => p + 1)}
            disabled={!hasMore}
            className="px-3 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded disabled:opacity-40 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
          >
            下一页
          </button>
        </div>
      )}
    </div>
  );
}
