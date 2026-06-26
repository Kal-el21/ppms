import axiosInstance from "../../../api/axiosInstance";
import type { ApiResponse, User } from "../../../types";

export const settingsApi = {
  updateProfile: async (payload: { full_name: string }) => {
    const res = await axiosInstance.put<ApiResponse<User>>("/me", payload);
    return res.data.data;
  },

  uploadProfilePhoto: async (file: File) => {
    const formData = new FormData();
    formData.append("photo", file);

    const res = await axiosInstance.post<ApiResponse<{ photo_url: string }>>("/me/photo", formData);
    return res.data.data;
  },

  toggle2FA: async (enabled: boolean) => {
    await axiosInstance.post("/me/toggle-2fa", { enabled });
  },

  toggleEmailNotification: async (enabled: boolean) => {
    await axiosInstance.post("/me/toggle-email-notification", { enabled });
  },
};
