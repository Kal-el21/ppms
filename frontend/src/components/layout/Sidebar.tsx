import { NavLink } from "react-router-dom";
import { useAuth } from "../../features/auth/context/AuthContext";
import type { ReactNode } from "react";
import UserMenu from "./UserMenu";

interface NavItem {
  label: string;
  path: string;
  roles: string[];
  icon: ReactNode;
  badge?: number;
}

interface NavGroup {
  label: string;
  items: NavItem[];
}

interface SidebarProps {
  onClose?: () => void;
}

const icon = (d: string) => (
  <svg className="h-4 w-4 flex-shrink-0 opacity-85" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
    <path d={d} />
  </svg>
);

const navGroups: NavGroup[] = [
  {
    label: "Overview",
    items: [
      {
        label: "Dashboard",
        path: "/dashboard",
        roles: ["ADMIN", "USER", "VIEWER"],
        icon: (
          <svg className="h-4 w-4 flex-shrink-0 opacity-85" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <rect x="3" y="3" width="7" height="9" rx="1.5" />
            <rect x="14" y="3" width="7" height="5" rx="1.5" />
            <rect x="14" y="12" width="7" height="9" rx="1.5" />
            <rect x="3" y="16" width="7" height="5" rx="1.5" />
          </svg>
        ),
      },
    ],
  },
  {
    label: "Workspace",
    items: [
      { label: "Divisions", path: "/divisions", roles: ["ADMIN", "USER", "VIEWER"], icon: icon("M3 7l4-4h6l4 4h4v13H3z") },
      { label: "Project Requests", path: "/project-requests", roles: ["ADMIN", "USER"], icon: icon("M3 7l4-4h6l4 4h4v13H3z") },
      { label: "Projects", path: "/projects", roles: ["ADMIN", "USER", "VIEWER"], icon: icon("M3 7l4-4h6l4 4h4v13H3z") },
      { label: "Deadlines", path: "/projects/deadline", roles: ["ADMIN", "USER", "VIEWER"], icon: icon("M12 8v4l3 3M12 2a10 10 0 100 20 10 10 0 000-20z") },
      { label: "Users", path: "/users", roles: ["ADMIN"], icon: icon("M16 21v-2a4 4 0 00-4-4H6a4 4 0 00-4 4v2M9 11a4 4 0 100-8 4 4 0 000 8z") },
    ],
  },
  {
    label: "Insights",
    items: [
      { label: "Reporting", path: "/reporting", roles: ["ADMIN", "USER"], icon: icon("M3 3v18h18M18 17V9M13 17V5M8 17v-3") },
      { label: "Approval Workflows", path: "/approval-workflows", roles: ["ADMIN"], icon: icon("M9 12l2 2 4-4M20 6L9 17l-5-5") },
      { label: "Audit Logs", path: "/audit-logs", roles: ["ADMIN"], icon: icon("M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z") },
    ],
  },
  {
    label: "Settings",
    items: [
      { label: "Notification Preferences", path: "/notification-preferences", roles: ["ADMIN", "USER", "VIEWER"], icon: icon("M18 8a6 6 0 00-12 0c0 7-3 9-3 9h18s-3-2-3-9M13.7 21a2 2 0 01-3.4 0") },
      { label: "Pagu Tahunan", path: "/settings/budget-years", roles: ["ADMIN"], icon: icon("M12 1v22M17 5H9.5a3.5 3.5 0 000 7h5a3.5 3.5 0 010 7H6") },
      { label: "Import / Export", path: "/settings/import-export", roles: ["ADMIN"], icon: icon("M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4M7 10l5 5 5-5M12 15V3") },
    ],
  },
];

export default function Sidebar({ onClose }: SidebarProps) {
  const { user } = useAuth();

  return (
    <aside className="w-[248px] flex-shrink-0 bg-surface border-r border-border flex flex-col h-screen sticky top-0">
      <div className="h-14 flex items-center justify-between gap-2.5 px-5 border-b border-border flex-shrink-0">
        <div className="flex items-center gap-2.5">
          <div className="h-[26px] w-[26px] rounded-[7px] flex items-center justify-center text-white font-bold text-[13px] flex-shrink-0 bg-gradient-to-br from-primary-600 to-danger-600">
            P
          </div>
          <span className="font-semibold text-[14.5px] tracking-tight">PPMS</span>
        </div>
        {onClose && (
          <button
            onClick={onClose}
            className="flex items-center justify-center h-7 w-7 rounded-md hover:bg-surface-secondary cursor-pointer"
            aria-label="Tutup sidebar"
          >
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M18 6L6 18M6 6l12 12" />
            </svg>
          </button>
        )}
      </div>

      <nav className="flex-1 overflow-y-auto py-4">
        {navGroups.map((group) => {
          const visibleItems = group.items.filter((item) => !user || item.roles.includes(user.system_role));
          if (visibleItems.length === 0) return null;

          return (
            <div key={group.label} className="px-3 mb-4">
              <div className="text-[11px] font-semibold uppercase tracking-wider text-ink-tertiary px-2 mb-1.5">
                {group.label}
              </div>
              {visibleItems.map((item) => (
                <NavLink
                  key={item.path}
                  to={item.path}
                  className={({ isActive }) =>
                    `relative flex items-center gap-2.5 px-2.5 py-[7px] mb-px rounded-md text-[13.5px] font-medium transition-colors ${
                      isActive
                        ? "bg-primary-50 text-primary-700 dark:bg-primary-900/30 dark:text-primary-400"
                        : "text-ink-secondary hover:bg-surface-secondary hover:text-ink-primary"
                    }`
                  }
                >
                  {({ isActive }) => (
                    <>
                      {isActive && (
                        <span className="absolute -left-3 top-1.5 bottom-1.5 w-[3px] bg-primary-600 rounded-r-[3px]" />
                      )}
                      {item.icon}
                      {item.label}
                      {item.badge && (
                        <span className="ml-auto text-[11px] font-semibold bg-danger-100 text-danger-700 dark:bg-danger-900/40 dark:text-danger-400 px-1.5 py-0.5 rounded-full">
                          {item.badge}
                        </span>
                      )}
                    </>
                  )}
                </NavLink>
              ))}
            </div>
          );
        })}
      </nav>

      <UserMenu />
    </aside>
  );
}
