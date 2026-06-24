import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { milestoneApi } from "../api/milestoneApi";
import { useToast } from "../../../components/ui/toast";

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