import axiosInstance from "../../../api/axiosInstance";
import type { ApiResponse } from "../../../types";
import type { Milestone } from "../types";

export const milestoneApi = {
  getByProject: async (projectId: number) => {
    const res = await axiosInstance.get<ApiResponse<Milestone[]>>(`/projects/${projectId}/milestones`);
    return res.data.data;
  },

  create: async (projectId: number, payload: { name: string; description: string }) => {
    const res = await axiosInstance.post<ApiResponse<Milestone>>(`/projects/${projectId}/milestones`, payload);
    return res.data.data;
  },

  changeStatus: async (projectId: number, milestoneId: number, status: string, version: number) => {
    await axiosInstance.patch(`/projects/${projectId}/milestones/${milestoneId}/status`, { status, version });
  },

  reorder: async (projectId: number, orderedIds: number[]) => {
    await axiosInstance.patch(`/projects/${projectId}/milestones/reorder`, { ordered_ids: orderedIds });
  },
};