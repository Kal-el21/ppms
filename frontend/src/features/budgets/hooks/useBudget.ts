import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { budgetApi } from "../api/budgetApi";

export function useBudget(projectId: number) {
  return useQuery({
    queryKey: ["projects", projectId, "budget"],
    queryFn: () => budgetApi.getByProject(projectId),
    enabled: !!projectId,
    retry: false, // budget mungkin belum dibuat (404 expected)
  });
}

export function useCreateBudget(projectId: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (allocatedBudget: number) => budgetApi.create(projectId, allocatedBudget),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["projects", projectId, "budget"] }),
  });
}

export function useTransactions(projectId: number, budgetId: number) {
  return useQuery({
    queryKey: ["projects", projectId, "budget", budgetId, "transactions"],
    queryFn: () => budgetApi.getTransactions(projectId, budgetId),
    enabled: !!projectId && !!budgetId,
  });
}

export function useCreateTransaction(projectId: number, budgetId: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      type: string;
      adjustment_type?: string;
      amount: number;
      reason?: string;
      description?: string;
      idempotency_key: string;
    }) => budgetApi.createTransaction(projectId, budgetId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, "budget"] });
    },
  });
}