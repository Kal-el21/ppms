import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { requestApi } from "../api/requestApi";
import { useToast } from "../../../components/ui/toast";
import type { CreateDraftPayload, UpdateDraftPayload, ReviewPayload, RevisePayload } from "../types";

export function useRequestList(page = 1, limit = 20, status?: string) {
  return useQuery({
    queryKey: ["project-requests", page, limit, status],
    queryFn: () => requestApi.getList(page, limit, status),
  });
}

export function useRequestDetail(id: number) {
  return useQuery({
    queryKey: ["project-requests", id],
    queryFn: () => requestApi.getById(id),
    enabled: !!id,
  });
}

export function useRevisionHistory(id: number) {
  return useQuery({
    queryKey: ["project-requests", id, "revisions"],
    queryFn: () => requestApi.getRevisions(id),
    enabled: !!id,
  });
}

export function useApprovalHistory(id: number) {
  return useQuery({
    queryKey: ["project-requests", id, "approvals"],
    queryFn: () => requestApi.getApprovals(id),
    enabled: !!id,
  });
}

export function useCreateDraft() {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: (payload: CreateDraftPayload) => requestApi.createDraft(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["project-requests"] });
      toast.success("Draft berhasil dibuat");
    },
    onError: (error: any) => {
      toast.error("Gagal membuat draft", error?.friendlyMessage);
    },
  });
}

export function useUpdateDraft(id: number) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: (payload: UpdateDraftPayload) => requestApi.updateDraft(id, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["project-requests"] });
      queryClient.invalidateQueries({ queryKey: ["project-requests", id] });
      toast.success("Draft berhasil disimpan");
    },
    onError: (error: any) => {
      toast.error("Gagal menyimpan draft", error?.friendlyMessage);
    },
  });
}

export function useSubmitRequest() {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: (id: number) => requestApi.submit(id),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: ["project-requests"] });
      queryClient.invalidateQueries({ queryKey: ["project-requests", id] });
      toast.success("Request berhasil disubmit", "Menunggu review dari Admin");
    },
    onError: (error: any) => {
      toast.error("Gagal submit request", error?.friendlyMessage);
    },
  });
}

export function useReviewRequest(id: number) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: (payload: ReviewPayload) => requestApi.review(id, payload),
    onSuccess: (_, payload) => {
      queryClient.invalidateQueries({ queryKey: ["project-requests"] });
      queryClient.invalidateQueries({ queryKey: ["project-requests", id] });
      queryClient.invalidateQueries({ queryKey: ["project-requests", id, "approvals"] });
      const label = payload.action === "APPROVED" ? "disetujui" : payload.action === "REJECTED" ? "ditolak" : "diminta revisi";
      toast.success(`Request berhasil ${label}`);
    },
    onError: (error: any) => {
      toast.error("Gagal submit review", error?.friendlyMessage);
    },
  });
}

export function useReviseRequest(id: number) {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: (payload: RevisePayload) => requestApi.revise(id, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["project-requests"] });
      queryClient.invalidateQueries({ queryKey: ["project-requests", id] });
      queryClient.invalidateQueries({ queryKey: ["project-requests", id, "revisions"] });
      toast.success("Revisi berhasil dikirim");
    },
    onError: (error: any) => {
      toast.error("Gagal mengirim revisi", error?.friendlyMessage);
    },
  });
}

export function useDeleteDraft() {
  const queryClient = useQueryClient();
  const toast = useToast();
  return useMutation({
    mutationFn: (id: number) => requestApi.deleteDraft(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["project-requests"] });
      toast.success("Draft berhasil dihapus");
    },
    onError: (error: any) => {
      toast.error("Gagal menghapus draft", error?.friendlyMessage);
    },
  });
}