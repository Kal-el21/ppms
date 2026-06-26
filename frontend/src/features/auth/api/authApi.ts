import axiosInstance from "../../../api/axiosInstance";
import type { ApiResponse } from "../../../types";
import type { LoginRequest, LoginResponse, LoginSuccessResponse } from "../types";

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
  verifyOtp: async (payload: { otp_session_token: string; otp_code: string }) => {
    const res = await axiosInstance.post<ApiResponse<LoginSuccessResponse>>("/auth/verify-otp", payload);
    return res.data.data;
  },

  resendOtp: async (otp_session_token: string) => {
    await axiosInstance.post("/auth/resend-otp", { otp_session_token });
  },
};