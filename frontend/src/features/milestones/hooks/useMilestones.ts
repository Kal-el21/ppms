import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { milestoneApi } from "../api/milestoneApi";

export function useMilestones(projectId: number) {
  return useQuery({
    queryKey: ["projects", projectId, "milestones"],
    queryFn: () => milestoneApi.getByProject(projectId),
    enabled: !!projectId,
  });
}

export function useCreateMilestone(projectId: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { name: string; description: string }) => milestoneApi.create(projectId, payload),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["projects", projectId, "milestones"] }),
  });
}