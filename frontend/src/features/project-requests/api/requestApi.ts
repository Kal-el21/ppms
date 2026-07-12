import axiosInstance from "../../../api/axiosInstance";
import type { ApiResponse, PaginatedResponse } from "../../../types";
import type {
  ProjectRequest,
  RequestRevision,
  RequestApproval,
  CreateDraftPayload,
  UpdateDraftPayload,
  ReviewPayload,
  RevisePayload,
} from "../types";

function toApiDate(value?: string | null) {
  return value ? new Date(value).toISOString() : null;
}

function normalizeRequestPayload<T extends CreateDraftPayload>(payload: T): T {
  return {
    ...payload,
    proposed_start_date: toApiDate(payload.proposed_start_date),
    proposed_end_date: toApiDate(payload.proposed_end_date),
  };
}

export const requestApi = {
  getList: async (page = 1, limit = 20, status?: string) => {
    const res = await axiosInstance.get<PaginatedResponse<ProjectRequest>>("/project-requests", {
      params: { page, limit, status },
    });
    return res.data;
  },

  getById: async (id: number) => {
    const res = await axiosInstance.get<ApiResponse<ProjectRequest>>(`/project-requests/${id}`);
    return res.data.data;
  },

  createDraft: async (payload: CreateDraftPayload) => {
    const res = await axiosInstance.post<ApiResponse<ProjectRequest>>(
      "/project-requests",
      normalizeRequestPayload(payload)
    );
    return res.data.data;
  },

  updateDraft: async (id: number, payload: UpdateDraftPayload) => {
    const res = await axiosInstance.put<ApiResponse<ProjectRequest>>(
      `/project-requests/${id}`,
      normalizeRequestPayload(payload)
    );
    return res.data.data;
  },

  submit: async (id: number) => {
    const res = await axiosInstance.post<ApiResponse<ProjectRequest>>(`/project-requests/${id}/submit`);
    return res.data.data;
  },

  review: async (id: number, payload: ReviewPayload) => {
    const res = await axiosInstance.post<ApiResponse<ProjectRequest>>(`/project-requests/${id}/review`, payload);
    return res.data.data;
  },

  revise: async (id: number, payload: RevisePayload) => {
    const res = await axiosInstance.post<ApiResponse<ProjectRequest>>(
      `/project-requests/${id}/revise`,
      normalizeRequestPayload(payload)
    );
    return res.data.data;
  },

  getRevisions: async (id: number) => {
    const res = await axiosInstance.get<ApiResponse<RequestRevision[]>>(`/project-requests/${id}/revisions`);
    return res.data.data;
  },

  getApprovals: async (id: number) => {
    const res = await axiosInstance.get<ApiResponse<RequestApproval[]>>(`/project-requests/${id}/approvals`);
    return res.data.data;
  },

  deleteDraft: async (id: number) => {
    await axiosInstance.delete(`/project-requests/${id}`);
  },
};
