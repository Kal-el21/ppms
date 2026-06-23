import axiosInstance from "../../../api/axiosInstance";
import type { ApiResponse } from "../../../types";
import type { LoginRequest, LoginResponse } from "../types";

export const authApi = {
  login: async (payload: LoginRequest) => {
    const res = await axiosInstance.post<ApiResponse<LoginResponse>>("/auth/login", payload);
    return res.data.data;
  },

  logout: async () => {
    // Tidak perlu kirim refresh_token lagi — backend membacanya dari cookie.
    await axiosInstance.post("/auth/logout");
  },

  changePassword: async (oldPassword: string, newPassword: string) => {
    await axiosInstance.post("/auth/change-password", {
      old_password: oldPassword,
      new_password: newPassword,
    });
  },
};