import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import { authApi } from "../api/authApi";
import { useAuth } from "../context/AuthContext";
import type { LoginRequest } from "../types";

export function useLogin() {
  const { login } = useAuth();
  const navigate = useNavigate();

  return useMutation({
    mutationFn: (payload: LoginRequest) => authApi.login(payload),
    onSuccess: (data) => {
      // Token sudah di-set sebagai httpOnly cookie oleh backend (Set-Cookie
      // header pada response ini) — tidak ada yang perlu disimpan manual
      // di sisi frontend selain data user untuk ditampilkan di UI.
      login(data.user);
      navigate("/dashboard");
    },
  });
}