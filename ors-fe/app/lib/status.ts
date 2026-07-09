import type { ReservationStatus } from "./api/reservations";

export type ServiceStatus = "active" | "inactive" | "pending" | "rejected";

export const STATUS_CONFIG: Record<string, { label: string; className: string }> = {
  active:    { label: "已上架",  className: "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-300" },
  inactive:  { label: "已下架",  className: "bg-gray-100 text-gray-500 dark:bg-gray-800 dark:text-gray-400" },
  pending:   { label: "待确认",  className: "bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-300" },
  confirmed: { label: "已确认",  className: "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-300" },
  completed: { label: "已完成",  className: "bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300" },
  cancelled: { label: "已取消",  className: "bg-gray-100 text-gray-500 dark:bg-gray-800 dark:text-gray-400" },
  rejected:  { label: "已拒绝",  className: "bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-300" },
};
