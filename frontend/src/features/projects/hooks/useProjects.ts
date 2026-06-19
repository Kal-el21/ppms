import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { projectApi } from "../api/projectApi";

export function useProjectList(page = 1, limit = 20, status?: string) {
  return useQuery({
    queryKey: ["projects", page, limit, status],
    queryFn: () => projectApi.getList(page, limit, status),
  });
}

export function useProjectDetail(id: number) {
  return useQuery({
    queryKey: ["projects", id],
    queryFn: () => projectApi.getById(id),
    enabled: !!id,
  });
}

export function useProjectMembers(id: number) {
  return useQuery({
    queryKey: ["projects", id, "members"],
    queryFn: () => projectApi.getMembers(id),
    enabled: !!id,
  });
}

export function useChangeProjectStatus(id: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ status, version }: { status: string; version: number }) =>
      projectApi.changeStatus(id, status, version),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["projects", id] }),
  });
}

export function useAddMember(id: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ userId, projectRole }: { userId: number; projectRole: string }) =>
      projectApi.addMember(id, userId, projectRole),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["projects", id, "members"] }),
  });
}

export function useRemoveMember(projectId: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (memberId: number) => projectApi.removeMember(projectId, memberId),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["projects", projectId, "members"] }),
  });
}