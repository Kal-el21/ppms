import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { divisionApi } from "../api/divisionApi";

export function useDivisions() {
  return useQuery({
    queryKey: ["divisions"],
    queryFn: divisionApi.getAll,
  });
}

export function useCreateDivision() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: divisionApi.create,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["divisions"] }),
  });
}