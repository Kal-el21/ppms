import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { taskApi } from "../api/taskApi";
import type { Task } from "../types";

export function useTasks(projectId: number) {
  return useQuery({
    queryKey: ["projects", projectId, "tasks"],
    queryFn: () => taskApi.getByProject(projectId),
    enabled: !!projectId,
  });
}

export function useCreateTask(projectId: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: Partial<Task>) => taskApi.create(projectId, payload),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["projects", projectId, "tasks"] }),
  });
}

export function useChangeTaskStatus(projectId: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ taskId, status, version }: { taskId: number; status: string; version: number }) =>
      taskApi.changeStatus(projectId, taskId, status, version),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["projects", projectId, "tasks"] }),
  });
}

export function useUpdateTaskProgress(projectId: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ taskId, progress, version }: { taskId: number; progress: number; version: number }) =>
      taskApi.updateProgress(projectId, taskId, progress, version),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["projects", projectId, "tasks"] }),
  });
}

export function useTaskComments(projectId: number, taskId: number, page = 1) {
  return useQuery({
    queryKey: ["projects", projectId, "tasks", taskId, "comments", page],
    queryFn: () => taskApi.getComments(projectId, taskId, page),
    enabled: !!taskId,
  });
}