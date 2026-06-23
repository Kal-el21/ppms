import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { projectApi } from "../api/projectApi";
import { useToast } from "../../../components/ui/toast";

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