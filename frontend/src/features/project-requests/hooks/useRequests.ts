import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { requestApi } from "../api/requestApi";
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
  return useMutation({
    mutationFn: (payload: CreateDraftPayload) => requestApi.createDraft(payload),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["project-requests"] }),
  });
}

export function useUpdateDraft(id: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: UpdateDraftPayload) => requestApi.updateDraft(id, payload),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["project-requests"] }),
  });
}

export function useSubmitRequest() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => requestApi.submit(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["project-requests"] }),
  });
}

export function useReviewRequest(id: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: ReviewPayload) => requestApi.review(id, payload),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["project-requests"] }),
  });
}

export function useReviseRequest(id: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: RevisePayload) => requestApi.revise(id, payload),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["project-requests"] }),
  });
}

export function useDeleteDraft() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => requestApi.deleteDraft(id),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["project-requests"] }),
  });
}