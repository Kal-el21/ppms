import { NavLink } from "react-router-dom";
import { useAuth } from "../../features/auth/context/AuthContext";
import type { ReactNode } from "react";
import {
  LayoutDashboard,
  Building2,
  FileText,
  FolderKanban,
  CalendarClock,
  Users,
  BarChart3,
  CheckCheck,
  ScrollText,
  Bell,
  Wallet,
  Download,
  X,
  PanelLeftClose,
  PanelLeftOpen,
} from "lucide-react";
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
  collapsed?: boolean;
  onToggleCollapse?: () => void;
}

const navGroups: NavGroup[] = [
  {
    label: "Overview",
    items: [
      {
        label: "Dashboard",
        path: "/dashboard",
        roles: ["ADMIN", "USER", "VIEWER"],
        icon: <LayoutDashboard className="h-4 w-4 flex-shrink-0" />,
      },
    ],
  },
  {
    label: "Workspace",
    items: [
      { label: "Divisions", path: "/divisions", roles: ["ADMIN", "USER", "VIEWER"], icon: <Building2 className="h-4 w-4 flex-shrink-0" /> },
      { label: "Project Requests", path: "/project-requests", roles: ["ADMIN", "USER"], icon: <FileText className="h-4 w-4 flex-shrink-0" /> },
      { label: "Projects", path: "/projects", roles: ["ADMIN", "USER", "VIEWER"], icon: <FolderKanban className="h-4 w-4 flex-shrink-0" /> },
      { label: "Deadlines", path: "/projects/deadline", roles: ["ADMIN", "USER", "VIEWER"], icon: <CalendarClock className="h-4 w-4 flex-shrink-0" /> },
      { label: "Users", path: "/users", roles: ["ADMIN"], icon: <Users className="h-4 w-4 flex-shrink-0" /> },
    ],
  },
  {
    label: "Insights",
    items: [
      { label: "Reporting", path: "/reporting", roles: ["ADMIN", "USER"], icon: <BarChart3 className="h-4 w-4 flex-shrink-0" /> },
      { label: "Approval Workflows", path: "/approval-workflows", roles: ["ADMIN"], icon: <CheckCheck className="h-4 w-4 flex-shrink-0" /> },
      { label: "Audit Logs", path: "/audit-logs", roles: ["ADMIN"], icon: <ScrollText className="h-4 w-4 flex-shrink-0" /> },
    ],
  },
  {
    label: "Settings",
    items: [
      { label: "Notification Preferences", path: "/notification-preferences", roles: ["ADMIN", "USER", "VIEWER"], icon: <Bell className="h-4 w-4 flex-shrink-0" /> },
      { label: "Pagu Tahunan", path: "/settings/budget-years", roles: ["ADMIN"], icon: <Wallet className="h-4 w-4 flex-shrink-0" /> },
      { label: "Import / Export", path: "/settings/import-export", roles: ["ADMIN"], icon: <Download className="h-4 w-4 flex-shrink-0" /> },
    ],
  },
];

export default function Sidebar({ onClose, collapsed = false, onToggleCollapse }: SidebarProps) {
  const { user } = useAuth();

  return (
    <aside
      className={`flex-shrink-0 bg-surface border-r border-border flex flex-col h-screen sticky top-0 transition-[width] duration-200 ${
        collapsed ? "w-[72px]" : "w-[248px]"
      }`}
    >
      <div className={`h-14 flex items-center justify-between border-b border-border flex-shrink-0 ${collapsed ? "px-2 justify-center" : "px-5"}`}>
        {collapsed ? (
          <span className="font-bold text-[13px] tracking-tight text-primary-700 dark:text-primary-400">IR</span>
        ) : (
          <div className="flex items-center gap-2.5">
            <img
              src="/assets/brand/logo-indonesia-re.png"
              alt="Indonesia Re"
              className="h-[28px] w-auto object-contain flex-shrink-0"
            />
            <span className="font-semibold text-[14.5px] tracking-tight">PPMS</span>
          </div>
        )}
        {!collapsed && onClose && (
          <button
            onClick={onClose}
            className="flex items-center justify-center h-7 w-7 rounded-md hover:bg-surface-secondary cursor-pointer text-ink-secondary"
            aria-label="Tutup sidebar"
          >
            <X className="h-3.5 w-3.5" />
          </button>
        )}
      </div>

      <nav className={`flex-1 overflow-y-auto py-4 ${collapsed ? "px-2" : ""}`}>
        {navGroups.map((group) => {
          const visibleItems = group.items.filter((item) => !user || item.roles.includes(user.system_role));
          if (visibleItems.length === 0) return null;

          return (
            <div key={group.label} className={collapsed ? "mb-3" : "px-3 mb-4"}>
              {!collapsed && (
                <div className="text-[11px] font-semibold uppercase tracking-wider text-ink-tertiary px-2.5 mb-1.5">
                  {group.label}
                </div>
              )}
              {visibleItems.map((item) => (
                <NavLink
                  key={item.path}
                  to={item.path}
                  title={collapsed ? item.label : undefined}
                  className={({ isActive }) =>
                    `relative flex items-center gap-2.5 rounded-md text-[13.5px] font-medium transition-colors ${
                      collapsed ? "justify-center px-0 py-2.5 mx-auto w-11" : "px-2.5 py-[7px] mb-px"
                    } ${
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
                      {!collapsed && item.label}
                      {!collapsed && item.badge && (
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

      <div className="border-t border-border flex-shrink-0">
        <button
          onClick={onToggleCollapse}
          className={`hidden lg:flex items-center gap-2 w-full text-ink-secondary hover:bg-surface-secondary hover:text-ink-primary transition-colors cursor-pointer ${
            collapsed ? "justify-center px-2 py-3" : "px-4 py-3 text-[12.5px] font-medium"
          }`}
          aria-label={collapsed ? "Perlebar sidebar" : "Ciutkan sidebar"}
          title={collapsed ? "Perlebar sidebar" : "Ciutkan sidebar"}
        >
          {collapsed ? <PanelLeftOpen className="h-4 w-4" /> : <PanelLeftClose className="h-4 w-4" />}
          {!collapsed && <span>Ciutkan</span>}
        </button>
        <UserMenu collapsed={collapsed} />
      </div>
    </aside>
  );
}
