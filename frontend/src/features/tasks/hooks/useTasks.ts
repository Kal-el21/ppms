import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { taskApi } from "../api/taskApi";
import type { TaskPriority } from "../types";
import { useToast } from "../../../components/ui/toast";

// Interface eksplisit untuk create task agar tidak bergantung ke Partial<Task>
// yang typing-nya kadang ketat di beberapa field.
export interface CreateTaskPayload {
  title: string;
  description?: string;
  priority: TaskPriority;           // "LOW" | "MEDIUM" | "HIGH" | "URGENT" — typed
  milestone_id?: number | null;
  start_date?: string | null;       // string ISO, bukan Date
  due_date?: string | null;         // string ISO, bukan Date
}

export function useTasks(projectId: number) {
  return useQuery({
    queryKey: ["projects", projectId, "tasks"],
    queryFn: () => taskApi.getByProject(projectId),
    enabled: !!projectId,
  });
}

export function useCreateTask(projectId: number) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: (payload: CreateTaskPayload) => taskApi.create(projectId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, "tasks"] });
      toast.success("Task berhasil ditambahkan");
    },
    onError: (error: any) => {
      toast.error("Gagal menambah task", error?.friendlyMessage);
    },
  });
}

export function useChangeTaskStatus(projectId: number) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: ({
      taskId,
      status,
      version,
    }: {
      taskId: number;
      status: string;
      version: number;
    }) => taskApi.changeStatus(projectId, taskId, status, version),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, "tasks"] });
    },
    onError: (error: any) => {
      toast.error("Gagal mengubah status task", error?.friendlyMessage);
    },
  });
}

export function useUpdateTaskProgress(projectId: number) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: ({
      taskId,
      progress,
      version,
    }: {
      taskId: number;
      progress: number;
      version: number;
    }) => taskApi.updateProgress(projectId, taskId, progress, version),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, "tasks"] });
    },
    onError: (error: any) => {
      toast.error("Gagal mengubah progress", error?.friendlyMessage);
    },
  });
}

export function useTaskComments(projectId: number, taskId: number, page = 1) {
  return useQuery({
    queryKey: ["projects", projectId, "tasks", taskId, "comments", page],
    queryFn: () => taskApi.getComments(projectId, taskId, page),
    enabled: !!taskId,
  });
}