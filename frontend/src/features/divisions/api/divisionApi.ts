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
};