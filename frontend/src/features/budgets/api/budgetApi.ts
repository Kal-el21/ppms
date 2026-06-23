import axiosInstance from "../../../api/axiosInstance";
import type { ApiResponse } from "../../../types";
import type { Budget, Transaction } from "../types";

export const budgetApi = {
  getByProject: async (projectId: number) => {
    const res = await axiosInstance.get<ApiResponse<Budget>>(`/projects/${projectId}/budget`);
    return res.data.data;
  },

  create: async (projectId: number, allocatedBudget: number) => {
    const res = await axiosInstance.post<ApiResponse<Budget>>(`/projects/${projectId}/budget`, {
      allocated_budget: allocatedBudget,
    });
    return res.data.data;
  },

  getTransactions: async (projectId: number, budgetId: number) => {
    const res = await axiosInstance.get<ApiResponse<Transaction[]>>(
      `/projects/${projectId}/budget/${budgetId}/transactions`
    );
    return res.data.data;
  },

  createTransaction: async (
    projectId: number,
    budgetId: number,
    payload: {
      type: string;
      adjustment_type?: string;
      amount: number;
      reason?: string;
      description?: string;
      idempotency_key: string;
    }
  ) => {
    const res = await axiosInstance.post<ApiResponse<Transaction>>(
      `/projects/${projectId}/budget/${budgetId}/transactions`,
      payload
    );
    return res.data.data;
  },
};