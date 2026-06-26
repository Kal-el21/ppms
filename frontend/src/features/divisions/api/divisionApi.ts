import axiosInstance from "../../../api/axiosInstance";
import type { ApiResponse, Division } from "../../../types";

export const divisionApi = {
  getAll: async () => {
    const res = await axiosInstance.get<ApiResponse<Division[]>>("/divisions");
    return res.data.data;
  },

  create: async (payload: { name: string; description: string }) => {
    const res = await axiosInstance.post<ApiResponse<Division>>("/divisions", payload);
    return res.data.data;
  },

  getById: async (id: number) => {
    const res = await axiosInstance.get<ApiResponse<Division>>(`/divisions/${id}`);
    return res.data.data;
  },

  update: async (id: number, payload: { name: string; description: string }) => {
    const res = await axiosInstance.put<ApiResponse<Division>>(`/divisions/${id}`, payload);
    return res.data.data;
  },

  delete: async (id: number) => {
    await axiosInstance.delete(`/divisions/${id}`);
  },
};