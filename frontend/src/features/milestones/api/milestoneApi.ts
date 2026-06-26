import axiosInstance from "../../../api/axiosInstance";
import type { ApiResponse } from "../../../types";
import type { Milestone } from "../types";
import type { CreateMilestonePayload } from "../hooks/useMilestones";

export const milestoneApi = {
  getByProject: async (projectId: number) => {
    const res = await axiosInstance.get<ApiResponse<Milestone[]>>(
      `/projects/${projectId}/milestones`
    );
    return res.data.data;
  },

  create: async (projectId: number, payload: CreateMilestonePayload) => {
    const body: any = { ...payload };
    if (body.start_date) body.start_date = new Date(body.start_date).toISOString();
    if (body.end_date) body.end_date = new Date(body.end_date).toISOString();
    const res = await axiosInstance.post<ApiResponse<Milestone>>(
      `/projects/${projectId}/milestones`,
      body
    );
    return res.data.data;
  },

  changeStatus: async (projectId: number, milestoneId: number, status: string, version: number) => {
    await axiosInstance.patch(
      `/projects/${projectId}/milestones/${milestoneId}/status`,
      { status, version }
    );
  },

  reorder: async (projectId: number, orderedIds: number[]) => {
    await axiosInstance.patch(
      `/projects/${projectId}/milestones/reorder`,
      { ordered_ids: orderedIds }
    );
  },

  update: async (projectId: number, milestoneId: number, payload: Partial<Milestone> & { version: number }) => {
    const body: any = { ...payload };
    if (body.start_date) body.start_date = new Date(body.start_date).toISOString();
    if (body.end_date) body.end_date = new Date(body.end_date).toISOString();
    const res = await axiosInstance.put<ApiResponse<Milestone>>(`/projects/${projectId}/milestones/${milestoneId}`, body);
    return res.data.data;
  },

  delete: async (projectId: number, milestoneId: number) => {
    await axiosInstance.delete(`/projects/${projectId}/milestones/${milestoneId}`);
  },
};
