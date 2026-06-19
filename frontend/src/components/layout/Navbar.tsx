import { useAuth } from "../../features/auth/context/AuthContext";
import { Button } from "../ui/button";
import { useNavigate } from "react-router-dom";
import { authApi } from "../../features/auth/api/authApi";
import NotificationBell from "./NotificationBell";
import GlobalSearchBar from "./GlobalSearchBar";

export default function Navbar() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = async () => {
    const refreshToken = localStorage.getItem("refresh_token");
    if (refreshToken) {
      try {
        await authApi.logout(refreshToken);
      } catch {
        // ignore error, proceed to clear local state anyway
      }
    }
    logout();
    navigate("/login");
  };

  return (
    <header className="flex h-16 items-center justify-between border-b bg-white px-6">
      <GlobalSearchBar />
      <NotificationBell />
      <div />
      <div className="flex items-center gap-4">
        <div className="text-right">
          <p className="text-sm font-medium">{user?.full_name}</p>
          <p className="text-xs text-slate-500">{user?.system_role}</p>
        </div>
        <Button variant="outline" size="sm" onClick={handleLogout}>
          Logout
        </Button>
      </div>
    </header>
  );
}