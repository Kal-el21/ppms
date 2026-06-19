import axiosInstance from "../../../api/axiosInstance";
import type { ApiResponse, PaginatedResponse } from "../../../types";
import type { Project, ProjectMember } from "../types";

export const projectApi = {
  getList: async (page = 1, limit = 20, status?: string) => {
    const res = await axiosInstance.get<PaginatedResponse<Project>>("/projects", {
      params: { page, limit, status },
    });
    return res.data;
  },

  getById: async (id: number) => {
    const res = await axiosInstance.get<ApiResponse<Project>>(`/projects/${id}`);
    return res.data.data;
  },

  update: async (id: number, payload: Partial<Project> & { version: number }) => {
    const res = await axiosInstance.put<ApiResponse<Project>>(`/projects/${id}`, payload);
    return res.data.data;
  },

  changeStatus: async (id: number, status: string, version: number) => {
    const res = await axiosInstance.patch<ApiResponse<Project>>(`/projects/${id}/status`, { status, version });
    return res.data.data;
  },

  getMembers: async (id: number) => {
    const res = await axiosInstance.get<ApiResponse<ProjectMember[]>>(`/projects/${id}/members`);
    return res.data.data;
  },

  
  addMember: async (id: number, userId: number, projectRole: string) => {
    const res = await axiosInstance.post<ApiResponse<ProjectMember>>(`/projects/${id}/members`, {
      user_id: userId,
      project_role: projectRole,
    });
    return res.data.data;
  },

  changeMemberRole: async (projectId: number, memberId: number, projectRole: string) => {
    await axiosInstance.patch(`/projects/${projectId}/members/${memberId}/role`, { project_role: projectRole });
  },

  removeMember: async (projectId: number, memberId: number) => {
    await axiosInstance.delete(`/projects/${projectId}/members/${memberId}`);
  },
};