import { NavLink } from "react-router-dom";
import { useAuth } from "../../features/auth/context/AuthContext";

const navItems = [
  { label: "Dashboard", path: "/dashboard", roles: ["ADMIN", "USER", "VIEWER"] },
  { label: "Users", path: "/users", roles: ["ADMIN"] },
  { label: "Divisions", path: "/divisions", roles: ["ADMIN", "USER", "VIEWER"] },
  { label: "Project Requests", path: "/project-requests", roles: ["ADMIN", "USER"] },
  { label: "Projects", path: "/projects", roles: ["ADMIN", "USER", "VIEWER"] },
  { label: "Dashboard", path: "/dashboard", roles: ["ADMIN", "USER", "VIEWER"] },
  { label: "Users", path: "/users", roles: ["ADMIN"] },
  { label: "Divisions", path: "/divisions", roles: ["ADMIN", "USER", "VIEWER"] },
  { label: "Project Requests", path: "/project-requests", roles: ["ADMIN", "USER"] },
  { label: "Projects", path: "/projects", roles: ["ADMIN", "USER", "VIEWER"] },
  { label: "Reporting", path: "/reporting", roles: ["ADMIN", "USER"] },
  { label: "Audit Logs", path: "/audit-logs", roles: ["ADMIN"] },
];

export default function Sidebar() {
  const { user } = useAuth();

  return (
    <aside className="w-64 border-r bg-white">
      <div className="flex h-16 items-center border-b px-6">
        <span className="text-lg font-bold">PPMS</span>
      </div>
      <nav className="flex flex-col gap-1 p-4">
        {navItems
          .filter((item) => !user || item.roles.includes(user.system_role))
          .map((item) => (
            <NavLink
              key={item.path}
              to={item.path}
              className={({ isActive }) =>
                `rounded-md px-3 py-2 text-sm font-medium transition-colors ${
                  isActive ? "bg-slate-900 text-white" : "text-slate-700 hover:bg-slate-100"
                }`
              }
            >
              {item.label}
            </NavLink>
          ))}
      </nav>
    </aside>
  );
}