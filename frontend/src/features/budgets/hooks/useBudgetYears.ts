import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { budgetYearApi } from "../api/budgetApi";
import { useToast } from "../../../components/ui/toast";

export function useBudgetYears() {
  return useQuery({
    queryKey: ["budget-years"],
    queryFn: () => budgetYearApi.getAll(),
  });
}

export function useCreateBudgetYear() {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: (payload: { year: number; capex_ceiling: number; opex_ceiling: number }) =>
      budgetYearApi.create(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["budget-years"] });
      toast.success("Pagu tahunan berhasil ditambahkan");
    },
    onError: (error: any) => {
      toast.error("Gagal menambah pagu", error?.friendlyMessage || error?.response?.data?.message);
    },
  });
}

export function useUpdateBudgetYear() {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: { capex_ceiling: number; opex_ceiling: number; version: number } }) =>
      budgetYearApi.update(id, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["budget-years"] });
      toast.success("Pagu tahunan berhasil diperbarui");
    },
    onError: (error: any) => {
      toast.error("Gagal memperbarui pagu", error?.friendlyMessage);
    },
  });
}

export function useDeleteBudgetYear() {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: (id: number) => budgetYearApi.remove(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["budget-years"] });
      toast.success("Pagu tahunan berhasil dihapus");
    },
    onError: (error: any) => {
      toast.error("Gagal menghapus pagu", error?.friendlyMessage);
    },
  });
}
