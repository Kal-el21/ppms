import axiosInstance from "../../../api/axiosInstance";
import type { ApiResponse } from "../../../types";
import type { Attachment, EntityType } from "../types";

export const attachmentApi = {
  upload: async (file: File, entityType: EntityType, entityId: number) => {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("entity_type", entityType);
    formData.append("entity_id", String(entityId));

    const res = await axiosInstance.post<ApiResponse<Attachment>>("/attachments", formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });
    return res.data.data;
  },

  getByEntity: async (entityType: EntityType, entityId: number) => {
    const res = await axiosInstance.get<ApiResponse<Attachment[]>>("/attachments", {
      params: { entity_type: entityType, entity_id: entityId },
    });
    return res.data.data;
  },

  getDownloadUrl: async (id: number) => {
    const res = await axiosInstance.get<ApiResponse<{ download_url: string; file_name: string; mime_type: string }>>(
      `/attachments/${id}/download`
    );
    return res.data.data;
  },

  delete: async (id: number) => {
    await axiosInstance.delete(`/attachments/${id}`);
  },
};