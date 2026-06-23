import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "../features/auth/context/AuthContext";
import type { SystemRole } from "../types";

interface ProtectedRouteProps {
  allowedRoles?: SystemRole[];
}

function PageLoader() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-surface-secondary">
      <div className="flex flex-col items-center gap-3">
        <div className="h-8 w-8 rounded-lg bg-gradient-to-br from-primary-600 to-danger-600 animate-pulse" />
        <p className="text-[13px] text-ink-tertiary">Memuat...</p>
      </div>
    </div>
  );
}

export default function ProtectedRoute({ allowedRoles }: ProtectedRouteProps) {
  const { isAuthenticated, isLoading, user } = useAuth();

  if (isLoading) return <PageLoader />;

  if (!isAuthenticated) return <Navigate to="/login" replace />;

  if (allowedRoles && user && !allowedRoles.includes(user.system_role)) {
    return <Navigate to="/dashboard" replace />;
  }

  return <Outlet />;
}