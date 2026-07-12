import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { projectApi } from "../api/projectApi";
import { useToast } from "../../../components/ui/toast";
import type { Project, CreateProjectDirectRequest } from "../types";

export interface ProjectListFilters {
  page?: number;
  limit?: number;
  search?: string;
  status?: string;
  budget_type?: string;
  initiation_type?: string;
  priority?: string;
  sort?: string;
  progress?: string;
}

export function useProjectList(filters: ProjectListFilters = {}) {
  const { page = 1, limit = 20, search, status, budget_type, initiation_type, priority, sort, progress } = filters;
  return useQuery({
    queryKey: ["projects", "list", page, limit, search, status, budget_type, initiation_type, priority, sort, progress],
    queryFn: () => projectApi.getList(filters),
  });
}

export function useProjectDeadlines(window = "90") {
  return useQuery({
    queryKey: ["projects", "deadlines", window],
    queryFn: () => projectApi.getDeadlines(window),
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
  const toast = useToast();
  return useMutation({
    mutationFn: ({ status, version }: { status: string; version: number }) =>
      projectApi.changeStatus(id, status, version),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects", id] });
      toast.success("Status project berhasil diubah");
    },
    onError: (error: any) => {
      toast.error("Gagal mengubah status", error?.friendlyMessage || error?.response?.data?.message);
    },
  });
}

export function useUpdateProject(id: number) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: (payload: Partial<Project> & { version: number }) =>
      projectApi.update(id, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects"] });
      queryClient.invalidateQueries({ queryKey: ["projects", id] });
      toast.success("Project berhasil diperbarui");
    },
    onError: (error: any) => {
      toast.error("Gagal memperbarui project", error?.friendlyMessage);
    },
  });
}

export function useAddMember(id: number) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: ({ userId, projectRole }: { userId: number; projectRole: string }) =>
      projectApi.addMember(id, userId, projectRole),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects", id, "members"] });
      toast.success("Member berhasil ditambahkan");
    },
    onError: (error: any) => {
      toast.error("Gagal menambah member", error?.friendlyMessage);
    },
  });
}

export function useRemoveMember(projectId: number) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: (memberId: number) => projectApi.removeMember(projectId, memberId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, "members"] });
      toast.success("Member berhasil dihapus");
    },
    onError: (error: any) => {
      toast.error("Gagal menghapus member", error?.friendlyMessage);
    },
  });
}

export function useChangeMemberRole(projectId: number) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: ({ memberId, projectRole }: { memberId: number; projectRole: string }) =>
      projectApi.changeMemberRole(projectId, memberId, projectRole),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, "members"] });
      toast.success("Role member berhasil diperbarui");
    },
    onError: (error: any) => {
      toast.error("Gagal memperbarui role member", error?.friendlyMessage);
    },
  });
}

export function useCreateDirectProject() {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: (payload: CreateProjectDirectRequest) => projectApi.createDirect(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects"] });
      toast.success("Project berhasil dibuat langsung");
    },
    onError: (error: any) => {
      toast.error("Gagal membuat project", error?.friendlyMessage || error?.response?.data?.message);
    },
  });
}

export function useDeleteProject() {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: (id: number) => projectApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects"] });
      toast.success("Project berhasil dihapus");
    },
    onError: (error: any) => {
      toast.error("Gagal menghapus project", error?.response?.data?.message || error?.message);
    },
  });
}
