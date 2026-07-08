import { Link } from "react-router";
import type { ReservationStatus } from "../api/reservations";

interface ReservationCardProps {
  id: number;
  serviceTitle: string;
  serviceId: number;
  providerName: string;
  startTime: string;
  endTime: string;
  status: ReservationStatus;
  note?: string;
  onCancel?: (id: number) => void;
  cancelPending?: boolean;
}

const STATUS_CONFIG: Record<ReservationStatus, { label: string; className: string }> = {
  pending:    { label: "待确认",  className: "bg-yellow-100 text-yellow-700" },
  confirmed:  { label: "已确认",  className: "bg-green-100 text-green-700" },
  completed:  { label: "已完成",  className: "bg-blue-100 text-blue-700" },
  cancelled:  { label: "已取消",  className: "bg-gray-100 text-gray-500" },
  rejected:   { label: "已拒绝",  className: "bg-red-100 text-red-700" },
};

function formatDate(iso: string) {
  const d = new Date(iso);
  const date = d.toLocaleDateString("zh-CN", { month: "long", day: "numeric", weekday: "short" });
  const time = d.toLocaleTimeString("zh-CN", { hour: "2-digit", minute: "2-digit" });
  return { date, time };
}

function isCancelable(status: ReservationStatus) {
  return status === "pending" || status === "confirmed";
}

export function ReservationCard({
  id,
  serviceTitle,
  serviceId,
  providerName,
  startTime,
  endTime,
  status,
  note,
  onCancel,
  cancelPending,
}: ReservationCardProps) {
  const cfg = STATUS_CONFIG[status];
  const start = formatDate(startTime);
  const end = formatDate(endTime);

  return (
    <div className="rounded-xl border border-gray-200 bg-white dark:border-gray-700 dark:bg-gray-800 hover:shadow-md transition-shadow overflow-hidden">
      <div className="p-4">
        <div className="flex items-start justify-between gap-3">
          <div className="min-w-0 flex-1">
            <Link
              to={`/services/${serviceId}`}
              className="font-semibold text-base text-gray-900 dark:text-gray-100 hover:text-blue-600 dark:hover:text-blue-400 leading-snug"
            >
              {serviceTitle}
            </Link>
            <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
              {providerName}
            </p>
          </div>
          <span className={`shrink-0 text-xs px-2 py-0.5 rounded ${cfg.className}`}>
            {cfg.label}
          </span>
        </div>

        <div className="mt-3 flex flex-wrap items-center gap-x-4 gap-y-1 text-sm text-gray-600 dark:text-gray-400">
          <span>
            {start.date} {start.time} — {end.time}
          </span>
        </div>

        {note && (
          <p className="mt-2 text-sm text-gray-500 dark:text-gray-400 line-clamp-2">
            <span className="text-gray-400">备注：</span>{note}
          </p>
        )}

        <div className="mt-3 flex items-center gap-2">
          {onCancel && isCancelable(status) && (
            <button
              onClick={() => onCancel(id)}
              disabled={cancelPending}
              className="text-xs text-red-500 border border-red-200 rounded px-2.5 py-1 hover:bg-red-50 dark:hover:bg-red-900/20 disabled:opacity-40 transition-colors"
            >
              取消预约
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
