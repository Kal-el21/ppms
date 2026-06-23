import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { userApi, type CreateUserPayload } from "../api/userApi";
import { useToast } from "../../../components/ui/toast";

export function useUsers(page = 1, limit = 20) {
  return useQuery({
    queryKey: ["users", page, limit],
    queryFn: () => userApi.getAll(page, limit),
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