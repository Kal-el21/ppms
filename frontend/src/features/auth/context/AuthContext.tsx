import { createContext, useContext, useState, useEffect, type ReactNode } from "react";
import axiosInstance from "../../../api/axiosInstance";
import type { User } from "../../../types";

interface AuthContextValue {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (user: User, csrfToken: string) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

// Catatan arsitektur: setelah migrasi ke httpOnly cookie, frontend TIDAK
// LAGI menyimpan access/refresh token sama sekali — baik di localStorage
// maupun di memory. Browser yang mengelola cookie sepenuhnya. Yang disimpan
// di sini hanya `user` (data non-sensitif untuk UI) yang di-cache di
// sessionStorage agar refresh halaman tidak langsung flash ke halaman
// login sebelum sempat verifikasi sesi ke backend.
export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Saat aplikasi pertama kali dimuat (hard refresh), tidak ada cara
    // membaca httpOnly cookie dari JS untuk tahu "apakah user masih login".
    // Solusinya: panggil endpoint yang butuh auth (mis. /auth/refresh atau
    // endpoint profile) dan biarkan backend yang menentukan validitas cookie.
    const restoreSession = async () => {
      const cachedUser = sessionStorage.getItem("ppms_user");
      if (cachedUser) {
        try {
          setUser(JSON.parse(cachedUser));
        } catch {
          sessionStorage.removeItem("ppms_user");
        }
      }

      try {
        const res = await axiosInstance.post("/auth/refresh");
        // refresh berhasil berarti cookie masih valid; backend sudah
        // menerbitkan cookie+csrf baru otomatis lewat Set-Cookie.
        if (res.data?.data?.csrf_token) {
          sessionStorage.setItem("ppms_csrf_ready", "1");
        }
      } catch {
        // Cookie tidak valid/sudah expired — bersihkan cache user lama
        setUser(null);
        sessionStorage.removeItem("ppms_user");
      } finally {
        setIsLoading(false);
      }
    };

    restoreSession();
  }, []);

  const login = (userData: User) => {
    setUser(userData);
    sessionStorage.setItem("ppms_user", JSON.stringify(userData));
  };

  const logout = () => {
    setUser(null);
    sessionStorage.removeItem("ppms_user");
    sessionStorage.removeItem("ppms_csrf_ready");
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated: !!user,
        isLoading,
        login,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return context;
}