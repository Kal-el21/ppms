import axiosInstance from "../../../api/axiosInstance";
import type { ApiResponse } from "../../../types";
import type { Handover } from "../types";

export const handoverApi = {
  getByProject: async (projectId: number, status?: string) => {
    const res = await axiosInstance.get<ApiResponse<Handover[]>>(`/projects/${projectId}/handovers`, {
      params: { status },
    });
    return res.data.data;
  },

  create: async (
    projectId: number,
    payload: { receiver_id?: number; description: string; delivery_date?: string; delivery_time?: string }
  ) => {
    const res = await axiosInstance.post<ApiResponse<Handover>>(`/projects/${projectId}/handovers`, payload);
    return res.data.data;
  },

  markReceived: async (projectId: number, handoverId: number, version: number) => {
    const res = await axiosInstance.patch<ApiResponse<Handover>>(
      `/projects/${projectId}/handovers/${handoverId}/received`,
      { version }
    );
    return res.data.data;
  },

  markReturned: async (projectId: number, handoverId: number, reason: string, version: number) => {
    const res = await axiosInstance.patch<ApiResponse<Handover>>(
      `/projects/${projectId}/handovers/${handoverId}/returned`,
      { reason, version }
    );
    return res.data.data;
  },
};