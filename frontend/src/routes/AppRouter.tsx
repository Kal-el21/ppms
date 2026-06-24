import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import ProtectedRoute from "./ProtectedRoute";
import DashboardLayout from "../components/layout/DashboardLayout";
import LoginPage from "../features/auth/pages/LoginPage";
import DashboardPage from "../features/dashboard/pages/DashboardPage";
import UsersPage from "../features/users/pages/UsersPage";
import DivisionsPage from "../features/divisions/pages/DivisionsPage";
import RequestListPage from "../features/project-requests/pages/RequestListPage";
import RequestFormPage from "../features/project-requests/pages/RequestFormPage";
import RequestDetailPage from "../features/project-requests/pages/RequestDetailPage";
import ProjectListPage from "../features/projects/pages/ProjectListPage";
import ProjectDetailPage from "../features/projects/pages/ProjectDetailPage";
import AuditLogPage from "@/features/audit/pages/AuditLogPage";
import ReportingPage from "@/features/reporting/pages/ReportingPage";
import NotificationPreferencesPage from "@/features/notifications/pages/NotificationPreferencesPage";
import SettingsPage from "@/features/settings/pages/SettingsPage";


function AppRouter() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />

        <Route element={<ProtectedRoute />}>
          <Route element={<DashboardLayout />}>
            <Route path="/dashboard" element={<DashboardPage />} />
            <Route path="/divisions" element={<DivisionsPage />} />
            <Route path="/project-requests" element={<RequestListPage />} />
            <Route path="/project-requests/new" element={<RequestFormPage />} />
            <Route path="/project-requests/:id" element={<RequestDetailPage />} />
            <Route path="/project-requests/:id/edit" element={<RequestFormPage />} />
            <Route path="/projects" element={<ProjectListPage />} />
            <Route path="/projects/:id" element={<ProjectDetailPage />} />
            <Route path="/reporting" element={<ReportingPage />} />
            <Route path="/notification-preferences" element={<NotificationPreferencesPage />} />
            <Route path="/settings" element={<SettingsPage />} />

            <Route element={<ProtectedRoute allowedRoles={["ADMIN"]} />}>
              <Route path="/users" element={<UsersPage />} />
              <Route path="/audit-logs" element={<AuditLogPage />} />
            </Route>
          </Route>
        </Route>

        <Route path="/" element={<Navigate to="/dashboard" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default AppRouter;