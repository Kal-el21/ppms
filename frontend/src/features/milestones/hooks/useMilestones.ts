import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { milestoneApi } from "../api/milestoneApi";
import { useToast } from "../../../components/ui/toast";
import type { Milestone } from "../types";

// Interface eksplisit agar form bisa kirim start_date / end_date opsional
export interface CreateMilestonePayload {
  name: string;
  description?: string;
  start_date?: string | null; // string ISO "YYYY-MM-DD", bukan Date object
  end_date?: string | null;
}

export function useMilestones(projectId: number) {
  return useQuery({
    queryKey: ["projects", projectId, "milestones"],
    queryFn: () => milestoneApi.getByProject(projectId),
    enabled: !!projectId,
  });
}

export function useCreateMilestone(projectId: number) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: (payload: CreateMilestonePayload) =>
      milestoneApi.create(projectId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, "milestones"] });
      toast.success("Milestone berhasil ditambahkan");
    },
    onError: (error: any) => {
      toast.error("Gagal menambah milestone", error?.friendlyMessage);
    },
  });
}

export function useChangeMilestoneStatus(projectId: number) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: ({
      milestoneId,
      status,
      version,
    }: {
      milestoneId: number;
      status: string;
      version: number;
    }) => milestoneApi.changeStatus(projectId, milestoneId, status, version),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, "milestones"] });
      toast.success("Status milestone berhasil diperbarui");
    },
    onError: (error: any) => {
      toast.error("Gagal memperbarui status milestone", error?.friendlyMessage);
    },
  });
}

export function useReorderMilestones(projectId: number) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: (orderedIds: number[]) => milestoneApi.reorder(projectId, orderedIds),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, "milestones"] });
      toast.success("Urutan milestone berhasil disimpan");
    },
    onError: (error: any) => {
      toast.error("Gagal menyimpan urutan milestone", error?.friendlyMessage);
    },
  });
}

export function useUpdateMilestone(projectId: number) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: ({
      milestoneId,
      payload,
    }: {
      milestoneId: number;
      payload: Partial<Milestone> & { version: number };
    }) => milestoneApi.update(projectId, milestoneId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, "milestones"] });
      toast.success("Milestone berhasil diperbarui");
    },
    onError: (error: any) => {
      toast.error("Gagal memperbarui milestone", error?.friendlyMessage);
    },
  });
}

export function useDeleteMilestone(projectId: number) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: (milestoneId: number) => milestoneApi.delete(projectId, milestoneId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, "milestones"] });
      toast.success("Milestone berhasil dihapus");
    },
    onError: (error: any) => {
      toast.error("Gagal menghapus milestone", error?.friendlyMessage);
    },
  });
}
