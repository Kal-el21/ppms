import axiosInstance from "../../../api/axiosInstance";
import type { ApiResponse } from "../../../types";
import type { Task } from "../types";
import type { CreateTaskPayload } from "../hooks/useTasks";

export const taskApi = {
  getByProject: async (projectId: number) => {
    const res = await axiosInstance.get<ApiResponse<Task[]>>(
      `/projects/${projectId}/tasks`
    );
    return res.data.data;
  },

  create: async (projectId: number, payload: CreateTaskPayload) => {
    const body: any = { ...payload };
    if (body.start_date) body.start_date = new Date(body.start_date).toISOString();
    if (body.due_date) body.due_date = new Date(body.due_date).toISOString();
    const res = await axiosInstance.post<ApiResponse<Task>>(
      `/projects/${projectId}/tasks`,
      body
    );
    return res.data.data;
  },

  changeStatus: async (
    projectId: number,
    taskId: number,
    status: string,
    version: number
  ) => {
    await axiosInstance.patch(
      `/projects/${projectId}/tasks/${taskId}/status`,
      { status, version }
    );
  },

  updateProgress: async (
    projectId: number,
    taskId: number,
    progress: number,
    version: number
  ) => {
    await axiosInstance.patch(
      `/projects/${projectId}/tasks/${taskId}/progress`,
      { progress, version }
    );
  },

  assignUsers: async (projectId: number, taskId: number, userIds: number[]) => {
    await axiosInstance.post(
      `/projects/${projectId}/tasks/${taskId}/assignees`,
      { user_ids: userIds }
    );
  },

  getComments: async (
    projectId: number,
    taskId: number,
    page = 1,
    limit = 20
  ) => {
    const res = await axiosInstance.get(
      `/projects/${projectId}/tasks/${taskId}/comments`,
      { params: { page, limit } }
    );
    return res.data as {
      data: any[];
      meta: { page: number; limit: number; total: number };
    };
  },

  addComment: async (projectId: number, taskId: number, comment: string) => {
    await axiosInstance.post(
      `/projects/${projectId}/tasks/${taskId}/comments`,
      { comment }
    );
  },

  getById: async (projectId: number, taskId: number) => {
    const res = await axiosInstance.get<ApiResponse<Task>>(`/projects/${projectId}/tasks/${taskId}`);
    return res.data.data;
  },

  update: async (projectId: number, taskId: number, payload: Partial<Task> & { version: number }) => {
    const body: any = { ...payload };
    if (body.start_date) body.start_date = new Date(body.start_date).toISOString();
    if (body.due_date) body.due_date = new Date(body.due_date).toISOString();
    const res = await axiosInstance.put<ApiResponse<Task>>(`/projects/${projectId}/tasks/${taskId}`, body);
    return res.data.data;
  },

  delete: async (projectId: number, taskId: number) => {
    await axiosInstance.delete(`/projects/${projectId}/tasks/${taskId}`);
  },
};
