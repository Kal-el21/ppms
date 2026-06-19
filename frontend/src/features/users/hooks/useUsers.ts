import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { userApi, type CreateUserPayload } from "../api/userApi";

export function useUsers(page = 1, limit = 20) {
  return useQuery({
    queryKey: ["users", page, limit],
    queryFn: () => userApi.getAll(page, limit),
  });
}

export function useCreateUser() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateUserPayload) => userApi.create(payload),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["users"] }),
  });
}

export function useDeactivateUser() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => userApi.deactivate(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["users"] }),
  });
}