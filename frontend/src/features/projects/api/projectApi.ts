import axiosInstance from "../../../api/axiosInstance";
import type { ApiResponse, PaginatedResponse } from "../../../types";
import type { Project, ProjectMember, ProjectDeadline, CreateProjectDirectRequest } from "../types";

export const projectApi = {
  getList: async (params: {
    page?: number;
    limit?: number;
    search?: string;
    status?: string;
    budget_type?: string;
    initiation_type?: string;
    priority?: string;
    sort?: string;
    progress?: string;
  } = {}) => {
    const res = await axiosInstance.get<PaginatedResponse<Project>>("/projects", {
      params,
    });
    return res.data;
  },

  getById: async (id: number) => {
    const res = await axiosInstance.get<ApiResponse<Project>>(`/projects/${id}`);
    return res.data.data;
  },

  getDeadlines: async (window = "90") => {
    const res = await axiosInstance.get<ApiResponse<ProjectDeadline[]>>("/projects/deadline", {
      params: { window },
    });
    return res.data.data;
  },

  update: async (id: number, payload: Partial<Project> & { version: number }) => {
    const body: any = { ...payload };
    if (body.start_date) body.start_date = new Date(body.start_date).toISOString();
    if (body.end_date) body.end_date = new Date(body.end_date).toISOString();
    const res = await axiosInstance.put<ApiResponse<Project>>(`/projects/${id}`, body);
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

  createDirect: async (payload: CreateProjectDirectRequest) => {
    const res = await axiosInstance.post<ApiResponse<Project>>("/admin/projects", payload);
    return res.data.data;
  },

  delete: async (id: number) => {
    await axiosInstance.delete(`/projects/${id}`);
  },
};
