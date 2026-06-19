import axiosInstance from "../../../api/axiosInstance";
import type { ApiResponse } from "../../../types";
import type { LoginRequest, LoginResponse, RefreshTokenResponse } from "../types";

export const authApi = {
  login: async (payload: LoginRequest) => {
    const res = await axiosInstance.post<ApiResponse<LoginResponse>>("/auth/login", payload);
    return res.data.data;
  },

  refreshToken: async (refreshToken: string) => {
    const res = await axiosInstance.post<ApiResponse<RefreshTokenResponse>>("/auth/refresh", {
      refresh_token: refreshToken,
    });
    return res.data.data;
  },

  logout: async (refreshToken: string) => {
    await axiosInstance.post("/auth/logout", { refresh_token: refreshToken });
  },

  changePassword: async (oldPassword: string, newPassword: string) => {
    await axiosInstance.post("/auth/change-password", {
      old_password: oldPassword,
      new_password: newPassword,
    });
  },
};