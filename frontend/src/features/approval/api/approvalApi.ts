import axiosInstance from "../../../api/axiosInstance";
import type { ApiResponse } from "../../../types";
import type { ApprovalLevel, ApprovalWorkflow } from "../types";

export const approvalApi = {
  getWorkflows: async () => {
    const res = await axiosInstance.get<ApiResponse<ApprovalWorkflow[]>>("/approval-workflows");
    return res.data.data;
  },

  createWorkflow: async (name: string) => {
    const res = await axiosInstance.post<ApiResponse<ApprovalWorkflow>>("/approval-workflows", { name });
    return res.data.data;
  },

  getLevels: async (workflowId: number) => {
    const res = await axiosInstance.get<ApiResponse<ApprovalLevel[]>>(`/approval-workflows/${workflowId}/levels`);
    return res.data.data;
  },

  createLevel: async (workflowId: number, payload: { level_order: number; role_required: string }) => {
    const res = await axiosInstance.post<ApiResponse<ApprovalLevel>>(
      `/approval-workflows/${workflowId}/levels`,
      payload
    );
    return res.data.data;
  },
};
