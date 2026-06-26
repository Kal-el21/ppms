import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { userApi, type CreateUserPayload } from "../api/userApi";
import { useToast } from "../../../components/ui/toast";

export function useUsers(page = 1, limit = 20, enabled = true) {
  return useQuery({
    queryKey: ["users", page, limit],
    queryFn: () => userApi.getAll(page, limit),
    enabled,
  });
}

export function useCreateUser() {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: (payload: CreateUserPayload) => userApi.create(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
      toast.success("User berhasil dibuat");
    },
    onError: (error: any) => {
      toast.error("Gagal membuat user", error?.friendlyMessage);
    },
  });
}

export function useDeactivateUser() {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: (id: number) => userApi.deactivate(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
      toast.success("User berhasil dinonaktifkan");
    },
    onError: (error: any) => {
      toast.error("Gagal menonaktifkan user", error?.friendlyMessage);
    },
  });
}

export function useUpdateUser() {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: { full_name: string; division_id?: number | null } }) =>
      userApi.update(id, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
      toast.success("User berhasil diperbarui");
    },
    onError: (error: any) => {
      toast.error("Gagal memperbarui user", error?.friendlyMessage);
    },
  });
}

export function useAssignUserRole() {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: ({ id, systemRole }: { id: number; systemRole: string }) =>
      userApi.assignRole(id, systemRole),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
      toast.success("Role user berhasil diperbarui");
    },
    onError: (error: any) => {
      toast.error("Gagal memperbarui role", error?.friendlyMessage);
    },
  });
}

export function useRestoreUser() {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: (id: number) => userApi.restore(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
      toast.success("User berhasil dipulihkan");
    },
    onError: (error: any) => {
      toast.error("Gagal memulihkan user", error?.friendlyMessage);
    },
  });
}
