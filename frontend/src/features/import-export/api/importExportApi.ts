import axiosInstance from "../../../api/axiosInstance";
import type { ApiResponse } from "../../../types";

export interface ImportResult {
  total_projects: number;
  imported: number;
  skipped: number;
  errors: string[];
  imported_project_ids: number[];
}

export const importExportApi = {
  // Download seluruh data project sebagai file JSON (ADMIN only).
  exportData: async () => {
    const res = await axiosInstance.get("/admin/export", { responseType: "blob" });
    return res.data as Blob;
  },

  // Upload file JSON hasil export untuk restore data (ADMIN only).
  importData: async (file: File) => {
    const formData = new FormData();
    formData.append("file", file);
    const res = await axiosInstance.post<ApiResponse<ImportResult>>("/admin/import", formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });
    return res.data.data;
  },
};
