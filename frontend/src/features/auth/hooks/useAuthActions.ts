import { useMutation } from "@tanstack/react-query";
import { authApi } from "../api/authApi";
import { useToast } from "../../../components/ui/toast";
import type { LoginRequest, LoginResponse } from "../types";

export function useLogin() {
  const toast = useToast();
  return useMutation<LoginResponse, any, LoginRequest>({
    mutationFn: (payload: LoginRequest) => authApi.login(payload),
    onError: (error: any) => {
      toast.error("Gagal login", error?.friendlyMessage || error?.message);
    },
  });
}

export function useVerifyOTP() {
  const toast = useToast();
  return useMutation({
    mutationFn: (payload: { otp_session_token: string; otp_code: string }) => authApi.verifyOtp(payload),
    onError: (error: any) => {
      toast.error("Gagal verifikasi OTP", error?.friendlyMessage || error?.message);
    },
  });
}

export function useResendOTP() {
  return useMutation({
    mutationFn: (otp_session_token: string) => authApi.resendOtp(otp_session_token),
    onError: () => {},
  });
}

export default useLogin;
