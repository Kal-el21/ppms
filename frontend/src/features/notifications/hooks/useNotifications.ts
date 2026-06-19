import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { notificationApi } from "../api/notificationApi";

export function useNotifications(page = 1, limit = 20, unreadOnly = false) {
  return useQuery({
    queryKey: ["notifications", page, limit, unreadOnly],
    queryFn: () => notificationApi.getList(page, limit, unreadOnly),
    refetchInterval: 30000, // poll setiap 30 detik (sederhana, tanpa WebSocket — sesuai roadmap v2.0)
  });
}

export function useMarkAsRead() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => notificationApi.markAsRead(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["notifications"] }),
  });
}

export function useMarkAllAsRead() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => notificationApi.markAllAsRead(),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["notifications"] }),
  });
}

export function useNotificationPreferences() {
  return useQuery({
    queryKey: ["notification-preferences"],
    queryFn: notificationApi.getPreferences,
  });
}

export function useUpdatePreference() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ type, enabled }: { type: string; enabled: boolean }) =>
      notificationApi.updatePreference(type, enabled),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["notification-preferences"] }),
  });
}