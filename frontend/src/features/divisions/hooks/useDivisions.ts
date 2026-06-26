import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { divisionApi } from "../api/divisionApi";

export function useDivisions() {
  return useQuery({
    queryKey: ["divisions"],
    queryFn: divisionApi.getAll,
  });
}

export function useDivision(id: number) {
  return useQuery({
    queryKey: ["divisions", id],
    queryFn: () => divisionApi.getById(id),
    enabled: !!id,
  });
}

export function useCreateDivision() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: divisionApi.create,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["divisions"] }),
  });
}

export function useUpdateDivision() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: { name: string; description: string } }) =>
      divisionApi.update(id, payload),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["divisions"] }),
  });
}

export function useDeleteDivision() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: divisionApi.delete,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["divisions"] }),
  });
}
