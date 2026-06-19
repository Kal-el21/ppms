import axiosInstance from "../../../api/axiosInstance";
import type { ApiResponse } from "../../../types";
import type { DashboardSummary } from "../types";

export const dashboardApi = {
  getSummary: async () => {
    const res = await axiosInstance.get<ApiResponse<DashboardSummary>>("/dashboard");
    return res.data.data;
  },
};