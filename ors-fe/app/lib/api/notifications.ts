import { request } from "./client";

export interface NotificationItem {
  id: number;
  user_id: number;
  title: string;
  content: string;
  type: string;
  is_read: boolean;
  created_at: string;
}

export interface UnreadCountResponse {
  count: number;
}

export function fetchNotifications(
  isRead?: boolean,
  page = 1,
  pageSize = 20
): Promise<NotificationItem[]> {
  const query = new URLSearchParams();
  if (isRead !== undefined) query.set("is_read", String(isRead));
  query.set("page", String(page));
  query.set("page_size", String(pageSize));
  return request<NotificationItem[]>(
    `/notifications?${query.toString()}`
  );
}

export function fetchUnreadCount(): Promise<UnreadCountResponse> {
  return request<UnreadCountResponse>("/notifications/unread-count");
}

export function markNotificationRead(
  id: number
): Promise<NotificationItem> {
  return request<NotificationItem>(`/notifications/${id}/read`, {
    method: "PUT",
  });
}

export function markAllNotificationsRead(): Promise<{ affected: number }> {
  return request<{ affected: number }>("/notifications/read-all", {
    method: "PUT",
  });
}
