import axiosInstance from "../../../api/axiosInstance";
import type { Notification, NotificationPreference } from "../types";

export const notificationApi = {
  getList: async (page = 1, limit = 20, unreadOnly = false) => {
    const res = await axiosInstance.get("/notifications", {
      params: { page, limit, unread_only: unreadOnly },
    });
    return res.data as {
      data: Notification[];
      meta: { page: number; limit: number; total: number; unread_count: number };
    };
  },

  markAsRead: async (id: number) => {
    await axiosInstance.patch(`/notifications/${id}/read`);
  },

  markAllAsRead: async () => {
    await axiosInstance.patch("/notifications/read-all");
  },

  getPreferences: async () => {
    const res = await axiosInstance.get<{ data: NotificationPreference[] }>("/notifications/preferences");
    return res.data.data;
  },

  updatePreference: async (type: string, enabled: boolean) => {
    await axiosInstance.put("/notifications/preferences", { type, enabled });
  },
};