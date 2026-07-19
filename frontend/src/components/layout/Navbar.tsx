import { useLocation, Link } from "react-router-dom";
import { useTheme } from "../../lib/theme";
import GlobalSearchBar from "./GlobalSearchBar";
import NotificationBell from "./NotificationBell";
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
  Settings,
  Wallet,
  Download,
  ChevronRight,
} from "lucide-react";

interface NavbarProps {
  onToggleSidebar: () => void;
  sidebarOpen: boolean;
}

const routeMap: Record<string, { label: string; icon?: typeof LayoutDashboard }> = {
  dashboard: { label: "Dashboard", icon: LayoutDashboard },
  divisions: { label: "Divisions", icon: Building2 },
  "project-requests": { label: "Project Requests", icon: FileText },
  projects: { label: "Projects", icon: FolderKanban },
  deadline: { label: "Deadlines", icon: CalendarClock },
  users: { label: "Users", icon: Users },
  reporting: { label: "Reporting", icon: BarChart3 },
  "approval-workflows": { label: "Approval Workflows", icon: CheckCheck },
  "audit-logs": { label: "Audit Logs", icon: ScrollText },
  "notification-preferences": { label: "Notifications", icon: Bell },
  settings: { label: "Settings", icon: Settings },
  "budget-years": { label: "Budget Years", icon: Wallet },
  "import-export": { label: "Import / Export", icon: Download },
  new: { label: "New", icon: undefined },
  edit: { label: "Edit", icon: undefined },
};

function getBreadcrumbs(pathname: string) {
  const parts = pathname.split("/").filter(Boolean);
  const crumbs: { label: string; to: string }[] = [{ label: "Home", to: "/dashboard" }];

  let accumulated = "";
  for (let i = 0; i < parts.length; i++) {
    const part = parts[i];
    if (/^\d+$/.test(part) || /^[0-9a-f]{8}-[0-9a-f]{4}/.test(part)) continue;
    accumulated += `/${part}`;
    const entry = routeMap[part];
    if (entry) {
      crumbs.push({ label: entry.label, to: accumulated });
    }
  }

  return crumbs;
}

export default function Navbar({ onToggleSidebar, sidebarOpen }: NavbarProps) {
  const { theme, toggleTheme } = useTheme();
  const location = useLocation();
  const breadcrumbs = getBreadcrumbs(location.pathname);

  return (
    <header className="h-14 flex-shrink-0 bg-surface border-b border-border flex items-center justify-between px-6 sticky top-0 z-10">
      <div className="flex items-center gap-3 min-w-0">
        <button
          onClick={onToggleSidebar}
          className="flex items-center justify-center h-8 w-8 rounded-md border border-border bg-surface hover:bg-surface-secondary cursor-pointer flex-shrink-0"
          aria-label={sidebarOpen ? "Tutup sidebar" : "Buka sidebar"}
        >
          {sidebarOpen ? (
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M18 6L6 18M6 6l12 12" />
            </svg>
          ) : (
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M3 12h18M3 6h18M3 18h18" />
            </svg>
          )}
        </button>

        <nav className="hidden md:flex items-center gap-1 text-[12.5px]">
          {breadcrumbs.map((crumb, i) => (
            <span key={crumb.to} className="flex items-center gap-1">
              {i > 0 && <ChevronRight className="h-3 w-3 text-ink-tertiary" />}
              {i === breadcrumbs.length - 1 ? (
                <span className="font-medium text-ink-primary">{crumb.label}</span>
              ) : (
                <Link to={crumb.to} className="text-ink-secondary hover:text-ink-primary transition-colors">
                  {crumb.label}
                </Link>
              )}
            </span>
          ))}
        </nav>

        <GlobalSearchBar />
      </div>

      <div className="flex items-center gap-2">
        <button
          onClick={toggleTheme}
          className="flex items-center gap-0.5 bg-surface-secondary border border-border rounded-full p-[3px] cursor-pointer"
          aria-label="Toggle theme"
        >
          <span
            className={`h-6 w-[26px] rounded-full flex items-center justify-center transition-colors ${
              theme === "light" ? "bg-surface text-primary-600 shadow-sm" : "text-ink-tertiary"
            }`}
          >
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2">
              <circle cx="12" cy="12" r="4" />
              <path d="M12 2v2M12 20v2M4.9 4.9l1.4 1.4M17.7 17.7l1.4 1.4M2 12h2M20 12h2M4.9 19.1l1.4-1.4M17.7 6.3l1.4-1.4" />
            </svg>
          </span>
          <span
            className={`h-6 w-[26px] rounded-full flex items-center justify-center transition-colors ${
              theme === "dark" ? "bg-surface text-primary-600 shadow-sm" : "text-ink-tertiary"
            }`}
          >
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2">
              <path d="M21 12.8A9 9 0 1111.2 3 7 7 0 0021 12.8z" />
            </svg>
          </span>
        </button>

        <NotificationBell />
      </div>
    </header>
  );
}
