import {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
  type ReactNode,
} from "react";
import axiosInstance from "../../../api/axiosInstance";
import type { User } from "../../../types";

interface AuthContextValue {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  // login sekarang hanya menerima User — token dikelola server via httpOnly cookie,
  // tidak ada yang perlu disimpan di frontend.
  login: (user: User) => void;
  // updateUser untuk update partial user data di context (setelah edit profile dll)
  // tanpa harus re-fetch dari API.
  updateUser: (partial: Partial<User>) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const cachedUser = sessionStorage.getItem("ppms_user");
    if (cachedUser) {
      try {
        setUser(JSON.parse(cachedUser));
      } catch {
        sessionStorage.removeItem("ppms_user");
      }
    }

    const restoreSession = async () => {
      try {
        await axiosInstance.post("/auth/refresh");
      } catch {
        setUser(null);
        sessionStorage.removeItem("ppms_user");
      } finally {
        setIsLoading(false);
      }
    };

    restoreSession();
  }, []);

  const login = useCallback((userData: User) => {
    setUser(userData);
    sessionStorage.setItem("ppms_user", JSON.stringify(userData));
  }, []);

  const updateUser = useCallback((partial: Partial<User>) => {
    setUser((prev) => {
      if (!prev) return prev;
      const updated = { ...prev, ...partial };
      sessionStorage.setItem("ppms_user", JSON.stringify(updated));
      return updated;
    });
  }, []);

  const logout = useCallback(() => {
    setUser(null);
    sessionStorage.removeItem("ppms_user");
  }, []);

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated: !!user,
        isLoading,
        login,
        updateUser,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) throw new Error("useAuth must be used within AuthProvider");
  return context;
}