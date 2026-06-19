import axiosInstance from "../../../api/axiosInstance";
import type { PaginatedResponse, User, ApiResponse } from "../../../types";

export interface CreateUserPayload {
  full_name: string;
  email: string;
  password: string;
  system_role: string;
  division_id?: number | null;
}

export const userApi = {
  getAll: async (page = 1, limit = 20) => {
    const res = await axiosInstance.get<PaginatedResponse<User>>("/users", {
      params: { page, limit },
    });
    return res.data;
  },

  create: async (payload: CreateUserPayload) => {
    const res = await axiosInstance.post<ApiResponse<User>>("/users", payload);
    return res.data.data;
  },

  assignRole: async (id: number, system_role: string) => {
    const res = await axiosInstance.patch<ApiResponse<User>>(`/users/${id}/role`, { system_role });
    return res.data.data;
  },

  deactivate: async (id: number) => {
    await axiosInstance.delete(`/users/${id}`);
  },

  restore: async (id: number) => {
    await axiosInstance.post(`/users/${id}/restore`);
  },
};