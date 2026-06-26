import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { approvalApi } from "../api/approvalApi";
import { useToast } from "../../../components/ui/toast";

export function useApprovalWorkflows() {
  return useQuery({
    queryKey: ["approval-workflows"],
    queryFn: approvalApi.getWorkflows,
  });
}

export function useCreateApprovalWorkflow() {
  const queryClient = useQueryClient();
  const toast = useToast();

  return useMutation({
    mutationFn: approvalApi.createWorkflow,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["approval-workflows"] });
      toast.success("Workflow berhasil dibuat");
    },
    onError: (error: any) => {
      toast.error("Gagal membuat workflow", error?.friendlyMessage);
    },
  });
}

export function useApprovalLevels(workflowId: number) {
  return useQuery({
    queryKey: ["approval-workflows", workflowId, "levels"],
    queryFn: () => approvalApi.getLevels(workflowId),
    enabled: !!workflowId,
  });
}

export function useCreateApprovalLevel(workflowId: number) {
  const queryClient = useQueryClient();
  const toast = useToast();

  return useMutation({
    mutationFn: (payload: { level_order: number; role_required: string }) =>
      approvalApi.createLevel(workflowId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["approval-workflows", workflowId, "levels"] });
      toast.success("Approval level berhasil ditambahkan");
    },
    onError: (error: any) => {
      toast.error("Gagal menambahkan level", error?.friendlyMessage);
    },
  });
}
