import axiosInstance from "../../../api/axiosInstance";
import type { ApiResponse } from "../../../types";
import type { Budget, Transaction, PortfolioBudgetYear, BudgetType } from "../types";

export const budgetApi = {
  getByProject: async (projectId: number) => {
    const res = await axiosInstance.get<ApiResponse<Budget>>(`/projects/${projectId}/budget`);
    return res.data.data;
  },

  create: async (
    projectId: number,
    payload: { allocated_budget: number; budget_type?: BudgetType; budget_name?: string }
  ) => {
    const res = await axiosInstance.post<ApiResponse<Budget>>(`/projects/${projectId}/budget`, payload);
    return res.data.data;
  },

  update: async (
    projectId: number,
    budgetId: number,
    payload: { allocated_budget: number; budget_type?: BudgetType; budget_name?: string; version: number }
  ) => {
    const res = await axiosInstance.put<ApiResponse<Budget>>(`/projects/${projectId}/budget/${budgetId}`, payload);
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

export const budgetYearApi = {
  getAll: async () => {
    const res = await axiosInstance.get<ApiResponse<PortfolioBudgetYear[]>>(`/admin/budget-years`);
    return res.data.data;
  },

  create: async (payload: { year: number; capex_ceiling: number; opex_ceiling: number }) => {
    const res = await axiosInstance.post<ApiResponse<PortfolioBudgetYear>>(`/admin/budget-years`, payload);
    return res.data.data;
  },

  update: async (
    id: number,
    payload: { capex_ceiling: number; opex_ceiling: number; version: number }
  ) => {
    const res = await axiosInstance.put<ApiResponse<PortfolioBudgetYear>>(`/admin/budget-years/${id}`, payload);
    return res.data.data;
  },

  remove: async (id: number) => {
    const res = await axiosInstance.delete<ApiResponse<null>>(`/admin/budget-years/${id}`);
    return res.data;
  },
};
