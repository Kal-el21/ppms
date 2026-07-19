import { useState, useRef, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../features/auth/context/AuthContext";
import { useTheme } from "../../lib/theme";
import { Avatar } from "../ui/avatar";
import { authApi } from "../../features/auth/api/authApi";

export default function UserMenu({ collapsed = false }: { collapsed?: boolean }) {
  const { user, logout } = useAuth();
  const { theme, toggleTheme } = useTheme();
  const navigate = useNavigate();
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const handleLogout = async () => {
    try { await authApi.logout(); } catch { /* ignore */ }
    logout();
    navigate("/login");
  };

  if (!user) return null;

  return (
    <div ref={ref} className={`relative ${collapsed ? "p-2" : "p-3"} border-t border-border`}>
      <button
        onClick={() => setOpen((v) => !v)}
        className={`w-full flex items-center gap-2.5 rounded-md p-1.5 hover:bg-surface-secondary transition-colors ${
          collapsed ? "justify-center" : ""
        }`}
        title={collapsed ? user.full_name : undefined}
      >
        {user.profile_photo_url ? (
          <img src={user.profile_photo_url} alt={user.full_name} className="h-7 w-7 rounded-full object-cover flex-shrink-0" />
        ) : (
          <Avatar name={user.full_name} size="sm" />
        )}
        {!collapsed && (
          <>
            <div className="flex-1 min-w-0 text-left">
              <p className="text-[13px] font-semibold truncate m-0">{user.full_name}</p>
              <p className="text-[11px] text-ink-tertiary m-0">{user.system_role}</p>
            </div>
            <svg
              width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"
              className={`text-ink-tertiary flex-shrink-0 transition-transform ${open ? "rotate-180" : ""}`}
            >
              <path d="M6 9l6 6 6-6" />
            </svg>
          </>
        )}
      </button>

      {open && (
        <div className={`absolute bottom-full mb-1.5 rounded-lg border border-border bg-surface shadow-lg overflow-hidden ${
          collapsed ? "left-2 right-2" : "left-3 right-3"
        }`}>
          <div className="px-3.5 py-3 border-b border-border">
            <p className="text-[12px] font-semibold m-0">{user.full_name}</p>
            <p className="text-[11px] text-ink-tertiary m-0 mt-0.5">{user.email}</p>
          </div>

          <div className="py-1">
            <button
              onClick={() => { navigate("/settings"); setOpen(false); }}
              className="w-full flex items-center gap-2.5 px-3.5 py-2 text-[13px] text-ink-primary hover:bg-surface-secondary transition-colors"
            >
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <circle cx="12" cy="8" r="4"/><path d="M4 20c0-4 3.6-7 8-7s8 3 8 7"/>
              </svg>
              Profile & Settings
            </button>

            <button
              onClick={() => { toggleTheme(); setOpen(false); }}
              className="w-full flex items-center gap-2.5 px-3.5 py-2 text-[13px] text-ink-primary hover:bg-surface-secondary transition-colors"
            >
              {theme === "light" ? (
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <path d="M21 12.8A9 9 0 1111.2 3 7 7 0 0021 12.8z"/>
                </svg>
              ) : (
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <circle cx="12" cy="12" r="4"/>
                  <path d="M12 2v2M12 20v2M4.9 4.9l1.4 1.4M17.7 17.7l1.4 1.4M2 12h2M20 12h2M4.9 19.1l1.4-1.4M17.7 6.3l1.4-1.4"/>
                </svg>
              )}
              Tampilan {theme === "light" ? "Gelap" : "Terang"}
            </button>
          </div>

          <div className="border-t border-border py-1">
            <button
              onClick={handleLogout}
              className="w-full flex items-center gap-2.5 px-3.5 py-2 text-[13px] text-danger-600 hover:bg-danger-50 dark:hover:bg-danger-900/20 transition-colors"
            >
              <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M9 21H5a2 2 0 01-2-2V5a2 2 0 012-2h4M16 17l5-5-5-5M21 12H9"/>
              </svg>
              Logout
            </button>
          </div>
        </div>
      )}
    </div>
  );
}