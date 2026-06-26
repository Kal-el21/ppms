import { useAuth } from "../../auth/context/AuthContext";
import { useDashboardSummary } from "../hooks/useDashboard";
import { Card, CardHeader, CardTitle, CardContent } from "../../../components/ui/card";
import { MetricCard } from "../../../components/ui/metric-card";
import { PageHeader } from "../../../components/shared/PageHeader";
import { Button } from "../../../components/ui/button";
import { MetricsSkeleton } from "@/components/ui/skeleton";

export default function DashboardPage() {
  const { user } = useAuth();
  const { data, isLoading } = useDashboardSummary();

  // Ganti semua bagian setelah skeleton loading:

  if (isLoading) {
    return (
      <div>
        <div className="mb-6">
          <div className="h-7 w-48 bg-surface-tertiary rounded-md animate-pulse mb-2" />
          <div className="h-4 w-64 bg-surface-tertiary rounded-md animate-pulse" />
        </div>
        <MetricsSkeleton />
      </div>
    );
  }

  // Tidak perlu "if (!data) return null" karena kita pakai optional chaining di bawah

  const totalProjects = data?.total_projects ?? 0;
  const activeProjects = data?.active_projects ?? 0;
  const completedProjects = data?.completed_projects ?? 0;
  const totalTasks = data?.total_tasks ?? 0;
  const completedTasks = data?.completed_tasks ?? 0;
  const pendingRequests = data?.pending_requests ?? 0;
  const overdueTasks = data?.overdue_tasks ?? 0;
  const budgetUsage = data?.total_budget_usage_percentage ?? 0;
  const recentActivities = data?.recent_activities ?? [];

  return (
    <div>
      <PageHeader
        title={`Selamat datang, ${user?.full_name?.split(" ")[0] ?? ""}` }
        subtitle="Ringkasan portofolio proyek Anda"
        actions={
          <>
            <Button variant="outline" size="sm">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4" />
                <path d="M7 10l5 5 5-5M12 15V3" />
              </svg>
              Export
            </Button>
            <Button variant="primary" size="sm">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M12 5v14M5 12h14" />
              </svg>
              Request baru
            </Button>
          </>
        }
      />

      <div className="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-7 gap-3.5 mb-7">
        <MetricCard label="Total Projects" value={totalProjects} iconColor="blue"
          icon={<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M3 7l4-4h6l4 4h4v13H3z"/></svg>}
        />
        <MetricCard label="Active Projects" value={activeProjects} iconColor="green"
          icon={<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><circle cx="12" cy="12" r="9"/><path d="M9 12l2 2 4-4"/></svg>}
        />
        <MetricCard label="Completed Projects" value={completedProjects} iconColor="teal"
          icon={<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M20 6L9 17l-5-5"/></svg>}
        />
        <MetricCard label="Total Tasks" value={totalTasks} iconColor="indigo"
          icon={<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><rect x="3" y="3" width="18" height="18" rx="2"/><path d="M9 12l2 2 4-4"/></svg>}
        />
        <MetricCard label="Completed Tasks" value={completedTasks} iconColor="emerald"
          icon={<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M22 11.08V12a10 10 0 1 1-5.3-8.8"/><path d="M22 4L12 14.01l-3-3"/></svg>}
        />
        <MetricCard label="Pending Requests" value={pendingRequests} iconColor="amber"
          delta={{ neutral: "Menunggu review" }}
          icon={<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 3"/></svg>}
        />
        <MetricCard label="Overdue Tasks" value={overdueTasks} iconColor="red"
          delta={overdueTasks > 0 ? { direction: "down", text: "Perlu perhatian" } : undefined}
          icon={<svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><circle cx="12" cy="12" r="9"/><path d="M12 8v4l2.5 2.5"/></svg>}
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-[1.4fr_1fr] gap-4">
        <Card>
          <CardHeader><CardTitle>Budget usage</CardTitle></CardHeader>
          <CardContent>
            <div className="flex items-baseline gap-2 mb-3">
              <span className="text-[28px] font-semibold tracking-tight">{budgetUsage.toFixed(1)}%</span>
              <span className="text-xs text-ink-tertiary">rata-rata seluruh project aktif</span>
            </div>
            <div className="h-2 rounded-full bg-surface-tertiary overflow-hidden">
              <div
                className="h-full rounded-full transition-all"
                style={{
                  width: `${Math.min(budgetUsage, 100)}%`,
                  background:
                    budgetUsage >= 100
                      ? "linear-gradient(90deg, #2563EB, #DC2626)"
                      : budgetUsage >= 80
                      ? "linear-gradient(90deg, #2563EB, #D97706)"
                      : "#2563EB",
                }}
              />
            </div>
          </CardContent>
        </Card>

        {recentActivities.length > 0 && (
          <Card>
            <CardHeader><CardTitle>Aktivitas terbaru</CardTitle></CardHeader>
            <CardContent>
              <div className="flex flex-col gap-3.5">
                {recentActivities.map((activity, idx) => (
                  <div key={idx} className="flex gap-2.5">
                    <span className="h-1.5 w-1.5 rounded-full bg-primary-600 mt-1.5 flex-shrink-0" />
                    <div className="min-w-0">
                      <p className="text-[12.5px] m-0">
                        <span className="font-semibold">{activity.module}</span> — {activity.action}
                      </p>
                      <p className="text-[11.5px] text-ink-tertiary mt-0.5 m-0">
                        {new Date(activity.created_at).toLocaleString("id-ID")}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}