import { useState, useMemo } from "react";
import { useQueries } from "@tanstack/react-query";
import {
  useProviderReservations,
  useConfirmReservation,
  useRejectReservation,
} from "../../../lib/hooks/use-provider-reservations";
import type { ReservationStatus } from "../../../lib/api/reservations";
import { fetchUserPublic } from "../../../lib/api/users";
import { STATUS_CONFIG } from "../../../lib/status";

const STATUS_FILTERS: Array<{ label: string; value: ReservationStatus | "" }> = [
  { label: "全部", value: "" },
  { label: "待确认", value: "pending" },
  { label: "已确认", value: "confirmed" },
  { label: "已完成", value: "completed" },
  { label: "已取消", value: "cancelled" },
  { label: "已拒绝", value: "rejected" },
];

export default function ProviderReservationsPage() {
  const [statusFilter, setStatusFilter] = useState<ReservationStatus | "">("");
  const [page, setPage] = useState(1);
  const pageSize = 10;

  const params = {
    ...(statusFilter ? { status: statusFilter as ReservationStatus } : {}),
    page,
    page_size: pageSize,
  };

  const { data, isLoading } = useProviderReservations(params);
  const confirmMutation = useConfirmReservation();
  const rejectMutation = useRejectReservation();

  const reservations = data?.items ?? [];
  const hasMore = (data?.items?.length ?? 0) >= pageSize;

  const reservationUserIds = useMemo(
    () => [...new Set(reservations.map((r) => r.user_id))],
    [reservations]
  );
  const userQueries = useQueries({
    queries: reservationUserIds.map((id) => ({
      queryKey: ["user-public", id],
      queryFn: () => fetchUserPublic(id),
    })),
  });
  const userMap = useMemo(() => {
    const map: Record<number, { name: string; avatar_url?: string }> = {};
    reservationUserIds.forEach((id, i) => {
      map[id] = userQueries[i]?.data ?? { name: `用户 ${id}` };
    });
    return map;
  }, [reservationUserIds, userQueries]);

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">用户预约</h1>

      <div className="flex items-center gap-4 mb-6">
        <select
          value={statusFilter}
          onChange={(e) => {
            setStatusFilter(e.target.value as ReservationStatus | "");
            setPage(1);
          }}
          className="border border-gray-300 dark:border-gray-600 rounded px-3 py-2 bg-white dark:bg-gray-800 text-sm"
        >
          {STATUS_FILTERS.map((f) => (
            <option key={f.value} value={f.value}>
              {f.label}
            </option>
          ))}
        </select>
      </div>

      {isLoading ? (
        <div className="flex justify-center py-20">
          <p className="text-gray-400 dark:text-gray-500">加载中...</p>
        </div>
      ) : reservations.length === 0 ? (
        <div className="flex justify-center py-20">
          <p className="text-gray-400 dark:text-gray-500">暂无预约</p>
        </div>
      ) : (
        <>
          <div className="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
            <table className="w-full text-sm">
              <thead>
                <tr className="bg-gray-50 dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
                  <th className="text-left px-4 py-3 font-medium text-gray-600 dark:text-gray-400">预约时间</th>
                  <th className="text-left px-4 py-3 font-medium text-gray-600 dark:text-gray-400">服务名称</th>
                  <th className="text-left px-4 py-3 font-medium text-gray-600 dark:text-gray-400">用户</th>
                  <th className="text-left px-4 py-3 font-medium text-gray-600 dark:text-gray-400">备注</th>
                  <th className="text-left px-4 py-3 font-medium text-gray-600 dark:text-gray-400">状态</th>
                  <th className="text-right px-4 py-3 font-medium text-gray-600 dark:text-gray-400">操作</th>
                </tr>
              </thead>
              <tbody>
                {reservations.map((r) => (
                  <tr
                    key={r.id}
                    className="border-b border-gray-100 dark:border-gray-800 hover:bg-gray-50 dark:hover:bg-gray-800/50"
                  >
                    <td className="px-4 py-3 whitespace-nowrap">
                      {new Date(r.start_time).toLocaleString("zh-CN")}
                    </td>
                    <td className="px-4 py-3 max-w-[200px] truncate">
                      {r.service?.title ?? "未知服务"}
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap">
                      {(() => {
                        const u = userMap[r.user_id] ?? { name: `用户 ${r.user_id}` };
                        return (
                          <span className="flex items-center gap-2">
                            <span className="w-6 h-6 rounded-full bg-gray-200 flex items-center justify-center text-xs text-gray-500 font-medium shrink-0 overflow-hidden">
                              {u.avatar_url ? (
                                <img src={u.avatar_url} alt={u.name} className="w-full h-full object-cover" />
                              ) : (
                                u.name[0]
                              )}
                            </span>
                            <span className="text-gray-900 dark:text-gray-100">{u.name}</span>
                          </span>
                        );
                      })()}
                    </td>
                    <td className="px-4 py-3 max-w-[150px] truncate text-gray-500 dark:text-gray-400">
                      {r.note || "-"}
                    </td>
                    <td className="px-4 py-3">
                      <span
                        className={`inline-block text-xs px-2 py-0.5 rounded ${
                          (STATUS_CONFIG[r.status] ?? STATUS_CONFIG.pending).className
                        }`}
                      >
                        {(STATUS_CONFIG[r.status] ?? STATUS_CONFIG.pending).label}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-right space-x-2">
                      {r.status === "pending" && (
                        <>
                          <button
                            onClick={() => confirmMutation.mutate(r.id)}
                            disabled={confirmMutation.isPending}
                            className="text-green-600 dark:text-green-400 hover:underline text-xs disabled:opacity-40"
                          >
                            确认
                          </button>
                          <button
                            onClick={() => rejectMutation.mutate(r.id)}
                            disabled={rejectMutation.isPending}
                            className="text-red-500 dark:text-red-400 hover:underline text-xs disabled:opacity-40"
                          >
                            拒绝
                          </button>
                        </>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {reservations.length > 0 && (
            <div className="flex items-center justify-center gap-2 mt-6">
              <button
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                disabled={page <= 1}
                className="px-3 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded disabled:opacity-40 hover:bg-gray-50 dark:hover:bg-gray-800"
              >
                上一页
              </button>
              <span className="text-sm text-gray-500 dark:text-gray-400">第 {page} 页</span>
              <button
                onClick={() => setPage((p) => p + 1)}
                disabled={!hasMore}
                className="px-3 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded disabled:opacity-40 hover:bg-gray-50 dark:hover:bg-gray-800"
              >
                下一页
              </button>
            </div>
          )}
        </>
      )}
    </div>
  );
}
