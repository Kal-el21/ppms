import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { handoverApi } from "../api/handoverApi";

export function useHandovers(projectId: number, status?: string) {
  return useQuery({
    queryKey: ["projects", projectId, "handovers", status],
    queryFn: () => handoverApi.getByProject(projectId, status),
    enabled: !!projectId,
  });
}

export function useCreateHandover(projectId: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { receiver_id?: number; description: string; delivery_date?: string; delivery_time?: string }) =>
      handoverApi.create(projectId, payload),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["projects", projectId, "handovers"] }),
  });
}

export function useMarkReceived(projectId: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ handoverId, version }: { handoverId: number; version: number }) =>
      handoverApi.markReceived(projectId, handoverId, version),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["projects", projectId, "handovers"] }),
  });
}