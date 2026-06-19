import axiosInstance from "../../../api/axiosInstance";

export const reportingApi = {
  // Laporan sistem-wide, ADMIN only
  generate: async (type: string, format: string) => {
    const res = await axiosInstance.post(
      "/reports/generate",
      { type, format },
      { responseType: "blob" }
    );
    return res.data as Blob;
  },

  // Laporan per-project, ADMIN atau PROJECT_MANAGER project tersebut
  generateForProject: async (projectId: number, type: string, format: string) => {
    const res = await axiosInstance.post(
      `/projects/${projectId}/reports/generate`,
      { type, format },
      { responseType: "blob" }
    );
    return res.data as Blob;
  },
};