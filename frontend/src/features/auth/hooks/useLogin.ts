import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import { authApi } from "../api/authApi";
import { useAuth } from "../context/AuthContext";
import type { LoginRequest, LoginResponse } from "../types";

export function useLogin() {
  const { login } = useAuth();
  const navigate = useNavigate();

  return useMutation<LoginResponse, any, LoginRequest>({
    mutationFn: (payload: LoginRequest) => authApi.login(payload),
    onSuccess: (data) => {
      if ("user" in data) {
        login(data.user);
        navigate("/dashboard");
      }
    },
  });
}
